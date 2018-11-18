package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Matcher interface {
	Match(left interface{}, right interface{}) (bool, error)
}

type MatcherEqual struct{}

func (m MatcherEqual) Match(left interface{}, right interface{}) (bool, error) {
	return left == right, nil
}

type MatcherLessThan struct{}

func (m MatcherLessThan) Match(left interface{}, right interface{}) (bool, error) {
	return left.(int64) < right.(int64), nil
}

type QueryItem struct {
	Key      string
	Operator Matcher
	Value    interface{}
}

type QueryTree struct {
	Negative bool
	MatchAll bool

	Item *QueryItem

	Left  *QueryTree
	Right *QueryTree
}

func (qt QueryTree) PrettyPrint() string {
	if qt.Item != nil {
		var str string
		str += qt.Item.Key
		switch qt.Item.Operator.(type) {
		case MatcherEqual:
			if qt.Negative {
				str += "<>"
			} else {
				str += "="
			}
		case MatcherLessThan:
			if qt.Negative {
				str += ">="
			} else {
				str += "<"
			}
		case nil:
		}

		switch vv := qt.Item.Value.(type) {
		case string:
			str += fmt.Sprintf("'%s' ", vv)
		case bool:
			str += fmt.Sprintf("%v ", vv)
		case byte:
			str += fmt.Sprintf("%v ", vv)
		case int64:
			str += fmt.Sprintf("%d ", vv)
		}

		return str

	}

	var match string
	if qt.MatchAll {
		match = "AND "
	} else {
		match = "OR "
	}

	var negative string
	if qt.Negative {
		negative = "!"
	}

	return negative + "(" + qt.Left.PrettyPrint() + match + qt.Right.PrettyPrint() + ") "
}

func (s *store) checkDuplicatedRecord(tableName string, primaryKey string, primaryKeyType ColumnTypes, primaryIndex int, primaryValue interface{}) ([]Record, error) {
	query := func(record Record) (bool, error) {
		desc, err := record.GetTableDesc()
		if err != nil {
			return false, err
		}

		if isDeletedItem(desc, record) {
			return false, nil
		}

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
		desc, err := record.GetTableDesc()
		if err != nil {
			return false, err
		}

		if isDeletedItem(desc, record) {
			return false, nil
		}

		return isQueryTreeMatch(desc, qt, record.Values)
	}
}

func isDeletedItem(desc *TableDesc, record Record) bool {
	for i, v := range record.Values {
		if desc.Columns[i].Name == "____flags____" {
			if v.(byte)&0x80 == 0x80 {
				return true
			}
		}
	}

	return false
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

	if qt.MatchAll {
		val := leftVal && rightVal
		if qt.Negative {
			return !val, nil
		}
		return val, nil
	}

	val := leftVal || rightVal
	if qt.Negative {
		return !val, nil
	}
	return val, nil
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
