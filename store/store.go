package store

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Store is store interface
type Store interface {
	RegisterTable(tableDesc TableDesc) error

	Insert(tableName string, record Record) (affectedRows int64, err error)

	Select(tableName string, qt *QueryTree) ([]Record, error)

	Update(tableName string, qt *QueryTree, setItems []UpdateSetItem) (affectedRows int64, err error)

	Delete(talbeName string, qt *QueryTree) (affectedRows int64, err error)
}

type store struct {
	InnerIDs map[string]int64

	logger *log.Logger
}

var (
	ErrDuplicatedRecord         = errors.New("duplicate record")
	ErrBoolNotSupportPrimaryKey = errors.New("bool type doesn't support primary key")
)

var tables []TableDesc

// NewStore creates new store implementation.
func NewStore(logger *log.Logger) Store {
	// TODO should recovery innerIDs from db file when db starts.
	return &store{
		InnerIDs: make(map[string]int64),
		logger:   logger,
	}
}

func (s *store) RegisterTable(tableDesc TableDesc) error {
	columns := make([]Column, 2)
	// flags
	columns[0] = Column{Name: "____flags____", Type: ColumnTypeByte}
	// inner id
	columns[1] = Column{Name: "____id____", Type: ColumnTypeInteger}

	columns = append(columns, tableDesc.Columns...)

	tableDesc.Columns = columns

	tables = append(tables, tableDesc)
	// TODO handle dumplicate register

	s.InnerIDs[tableDesc.Name] = 0
	return nil
}

// Insert
// 1. the order of values must be same with the order of table desc
func (s *store) Insert(tableName string, record Record) (affectedRows int64, err error) {
	innerValues := make([]interface{}, 2)
	innerValues[0] = byte(0)
	innerValues[1] = s.InnerIDs[tableName] + 1
	record.Values = append(innerValues, record.Values...)

	a, e := s.insert(record)
	if e != nil {
		return 0, e
	}

	s.InnerIDs[tableName]++
	return a, nil
}

func (s *store) insert(record Record) (affectedRows int64, err error) {
	tableDesc, err := record.GetTableDesc()
	if err != nil {
		return 0, err
	}

	primaryKey, ptype, columnIndex := tableDesc.GetPrimaryKey()
	if primaryKey != "" {
		if _, err := s.checkDuplicatedRecord(record.TableName, primaryKey, ptype, columnIndex, record.Values[columnIndex]); err != nil {
			return 0, err
		}
	}

	recordBytes, err := getRecordBytes(record)
	if err != nil {
		return 0, err
	}

	tableFile, err := tableDesc.GetTableFile()
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(tableFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	_, err = f.Write(recordBytes)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (s *store) checkDuplicatedRecord(tableName string, primaryKey string, primaryKeyType ColumnTypes, primaryIndex int, primaryValue interface{}) ([]Record, error) {
	query := func(record Record) (bool, error) {
		return record.Values[primaryIndex] == primaryValue, nil
	}

	records, err := s.scanRecords(tableName, query, nil)
	if err != nil {
		return nil, err
	}

	if len(records) > 0 {
		return records, ErrDuplicatedRecord
	}

	return nil, nil
}

func query(qt *QueryTree) hitTarget {
	return func(record Record) (bool, error) {
		desc, _ := record.GetTableDesc()

		for i, v := range record.Values {
			if desc.Columns[i].Name == "____flags____" {
				if v.(byte)&0x80 == 0x80 {
					return false, nil
				}
			}
		}

		return isQueryTreeMatch(desc, qt, record.Values)
	}
}

func isQueryTreeMatch(tableDesc *TableDesc, qt *QueryTree, recordValues []interface{}) (bool, error) {
	if qt == nil {
		return true, nil
	}

	if qt.Item != nil {
		qItem := qt.Item

		index, err := tableDesc.IndexOfColumn(qItem.Key)
		if err != nil {
			return false, err
		}

		field := recordValues[index]
		match, err := qt.Item.Operator.Match(qItem.Value, field)
		if err != nil {
			return false, err
		}

		if qt.Negative {
			return !match, nil
		}

		return match, nil
	}

	leftVal, err := isQueryTreeMatch(tableDesc, qt.Left, recordValues)
	if err != nil {
		return false, err
	}
	rightVal, err := isQueryTreeMatch(tableDesc, qt.Right, recordValues)
	if err != nil {
		return false, err
	}

	if qt.Left.Negative {
		leftVal = !leftVal
	}
	if qt.Right.Negative {
		rightVal = !rightVal
	}

	if qt.MatchAll {
		return leftVal && rightVal, nil
	}
	return leftVal || rightVal, nil
}

func (s *store) Select(tableName string, qt *QueryTree) ([]Record, error) {
	records, err := s.scanRecords(tableName, query(qt), nil)
	if err != nil {
		return nil, err
	}

	for i := range records {
		outValues := records[i].Values[2:]
		records[i].Values = outValues
	}

	return records, nil
}

func (s *store) Update(tableName string, qt *QueryTree, setItems []UpdateSetItem) (affectedRows int64, err error) {
	return s.update(tableName, qt, setItems)
}

func (s *store) update(tableName string, qt *QueryTree, setItems []UpdateSetItem) (affectedRows int64, err error) {
	tableDesc, err := GetTableDescFromTableName(tableName)
	if err != nil {
		return 0, err
	}

	updateReplacer := func(record Record) (newRecord *Record, err error) {
		newRecord = &Record{
			TableName: tableName,
			Values:    make([]interface{}, len(record.Values)),
		}

		for i, c := range tableDesc.Columns {
			for _, si := range setItems {
				if si.Name == c.Name {
					outRecord := getOutRecord(record)
					newValue, err := si.Value(outRecord)
					if err != nil {
						return nil, err
					}

					primaryKeyName, pType, index := tableDesc.GetPrimaryKey()

					var containPrimaryKeyColumnUpdate bool
					for _, si := range setItems {
						if si.Name == primaryKeyName {
							containPrimaryKeyColumnUpdate = true
						}
					}

					if containPrimaryKeyColumnUpdate {
						duplicatedRecords, err := s.checkDuplicatedRecord(tableName, primaryKeyName, pType, index, newValue)
						if err != nil {
							if err == ErrDuplicatedRecord {
								if len(duplicatedRecords) == 1 && duplicatedRecords[0].Values[1] == record.Values[1] {
									goto Normal
								}
							}
							return nil, err
						}
					}

				Normal:
					newRecord.Values[i] = newValue
				} else {
					newRecord.Values[i] = record.Values[i]
				}
			}
		}

		return newRecord, err
	}

	records, err := s.scanRecords(tableName, query(qt), updateReplacer)
	if err != nil {
		return 0, err
	}

	return int64(len(records)), nil

}

func (s *store) Delete(tableName string, qt *QueryTree) (affectedRows int64, err error) {
	deleteItemFn := func(r Record) (interface{}, error) {
		return byte(0x80), nil
	}

	return s.update(tableName, qt, []UpdateSetItem{
		UpdateSetItem{
			Name:  "____flags____",
			Value: deleteItemFn,
		},
	})
}

type hitTarget func(record Record) (bool, error)
type replacer func(record Record) (newRecord *Record, err error)

func (s *store) scanRecords(tableName string, ht hitTarget, r replacer) ([]Record, error) {
	tableDesc, err := GetTableDescFromTableName(tableName)
	if err != nil {
		return nil, err
	}

	tableFile, err := tableDesc.GetTableFile()
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(tableFile, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	totalBytes := tableDesc.GetTotalBytes()

	bs := make([]byte, totalBytes)
	offset := int64(0)

	records := make([]Record, 0)
	for {
		_, err := f.ReadAt(bs, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		record := Record{
			TableName: tableName,
			Values:    make([]interface{}, 0),
		}

		for _, c := range tableDesc.Columns {
			columnOffset, err := tableDesc.OffsetOfColumn(c.Name)
			if err != nil {
				return nil, err
			}

			columnBytes := bs[columnOffset:(columnOffset + sizeOf(c.Type))]

			switch c.Type {
			case ColumnTypeString:
				var str string
				index := bytes.IndexByte(columnBytes, byte(0))
				if index != -1 {
					str = string(columnBytes[:index])
				} else {
					str = string(columnBytes)
				}
				record.Values = append(record.Values, str)
			case ColumnTypeInteger:
				buf := bytes.NewBuffer(columnBytes)
				var x int64
				binary.Read(buf, binary.BigEndian, &x)
				record.Values = append(record.Values, x)
			case ColumnTypeByte:
				record.Values = append(record.Values, columnBytes[0])
			case ColumnTypeBool:
				v := columnBytes[0] != byte(0)
				record.Values = append(record.Values, v)
			}
		}

		hit, err := ht(record)
		if err != nil {
			return nil, err
		}

		if hit {
			if r != nil {
				newRecord, err := r(record)

				if err != nil {
					return nil, err
				}

				recordBytes, err := getRecordBytes(*newRecord)
				if err != nil {
					return nil, err
				}

				_, err = f.WriteAt(recordBytes, offset)
				if err != nil {
					return nil, err
				}
			}

			records = append(records, record)
		}

		offset += int64(totalBytes)
	}

	return records, nil
}

func getRecordBytes(record Record) ([]byte, error) {
	bs := make([]byte, 0)

	for _, v := range record.Values {
		switch vv := v.(type) {
		case string:
			b255 := make([]byte, 0)
			bvv := []byte(vv)
			vvLen := len(bvv)
			b255 = append(b255, bvv...)
			b255 = append(b255, make([]byte, 255-vvLen)...)
			bs = append(bs, b255...)
		case int64:
			buf := bytes.NewBuffer([]byte{})
			binary.Write(buf, binary.BigEndian, vv)
			bs = append(bs, buf.Bytes()...)
		case byte:
			bs = append(bs, vv)
		case bool:
			var b byte
			if vv {
				b = byte(1)
			} else {
				b = byte(0)
			}
			bs = append(bs, b)
		default:
			return nil, errors.New("invalid type")
		}
	}

	return bs, nil
}

func getWorkingPath() (path string, err error) {
	ePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ePath), nil
}

func getOutRecord(record Record) Record {
	newRecord := record
	newRecord.Values = record.Values[2:]

	return newRecord
}
