package store

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
	"path/filepath"
)

// Store is store interface
type Storage interface {
	RegisterTable(tableDesc TableDesc) error

	Insert(tableName string, record Record) (affectedRows int64, err error)

	Select(tableName string, qt *QueryTree) ([]Record, error) // TODO support fields

	Update(tableName string, qt *QueryTree, setItems []SetItem) (affectedRows int64, err error)

	Delete(talbeName string, qt *QueryTree) (affectedRows int64, err error)
}

type Store struct {
	logger *log.Logger
}

var (
	ErrDuplicatedRecord         = errors.New("duplicate record")
	ErrBoolNotSupportPrimaryKey = errors.New("bool type doesn't support primary key")
	ErrDuplicatedTable          = errors.New("duplicate table name")
)

// NewStore creates new store implementation.
func NewStore(logger *log.Logger) *Store {
	return &Store{
		logger: logger,
	}
}

func (s *Store) RegisterTable(tableDesc TableDesc) error {
	for _, t := range tables {
		if t.Name == tableDesc.Name {
			return ErrDuplicatedTable
		}
	}

	columns := make([]Column, 2)
	// flags
	columns[0] = Column{Name: "____flags____", Type: ColumnTypeByte}
	// inner id
	columns[1] = Column{Name: "____id____", Type: ColumnTypeInteger}

	columns = append(columns, tableDesc.Columns...)
	tableDesc.Columns = columns
	tables = append(tables, tableDesc)

	err := SaveMetadata()
	if err != nil {
		return err
	}

	return nil
}

// Insert
// 1. the order of values must be same with the order of table desc
func (s *Store) Insert(tableName string, record Record) (affectedRows int64, err error) {
	tableDesc, err := record.GetTableDesc()
	if err != nil {
		return 0, err
	}

	innerValues := make([]interface{}, 2)
	innerValues[0] = byte(0)
	innerValues[1] = tableDesc.MaxInnerID + 1
	record.Values = append(innerValues, record.Values...)

	a, e := s.insert(record)
	if e != nil {
		return 0, e
	}

	tableDesc.MaxInnerID++
	return a, nil
}

func (s *Store) insert(record Record) (affectedRows int64, err error) {
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

func (s *Store) Update(tableName string, qt *QueryTree, setItems []SetItem) (affectedRows int64, err error) {
	return s.update(tableName, qt, setItems)
}

func (s *Store) update(tableName string, qt *QueryTree, setItems []SetItem) (affectedRows int64, err error) {
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

func (s *Store) Delete(tableName string, qt *QueryTree) (affectedRows int64, err error) {
	deleteItemFn := func(r Record) (interface{}, error) {
		return byte(0x80), nil
	}

	return s.update(tableName, qt, []SetItem{
		SetItem{
			Name:  "____flags____",
			Value: deleteItemFn,
		},
	})
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
