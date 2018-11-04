package store

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Store is store interface
type Store interface {
	RegisterTable(tableDesc TableDesc) error

	Insert(record Record) (savedRecord *Record, err error)
}

type store struct{}

var (
	ErrDuplicatedRecord         = errors.New("duplicate record")
	ErrBoolNotSupportPrimaryKey = errors.New("bool type doesn't support primary key")
)

type ColumnTypes int

const (
	_ ColumnTypes = iota
	ColumnTypeString
	ColumnTypeBool
	ColumnTypeInteger
)

type Column struct {
	Name       string
	Type       ColumnTypes
	PrimaryKey bool
}

type TableDesc struct {
	Name    string
	Columns []Column
}

func (t TableDesc) GetPrimaryKey() (name string, ptype ColumnTypes, index int) {
	for i, c := range t.Columns {
		if c.PrimaryKey {
			return c.Name, c.Type, i
		}
	}

	return "", ColumnTypeBool, 0
}

func (t TableDesc) OffsetOfColumn(columnName string) (int, error) {
	var offset int
	for _, c := range t.Columns {
		if c.Name == columnName {
			return offset, nil
		}

		offset += sizeOf(c.Type)
	}

	return -1, fmt.Errorf("unkonwn column %s", columnName)
}

func (t TableDesc) GetTotalBytes() int {
	var total int
	for _, c := range t.Columns {
		total += sizeOf(c.Type)
	}

	return total
}

func sizeOf(t ColumnTypes) int {
	switch t {
	case ColumnTypeString:
		return 255
	case ColumnTypeInteger:
		return 8
	case ColumnTypeBool:
		return 1
	}

	panic("unsupport type")
}

func (t TableDesc) GetTableFile() (string, error) {
	workingPath, err := getWorkingPath()
	if err != nil {
		return "", err
	}

	tableFile := filepath.Join(workingPath, t.Name)
	return tableFile, nil
}

type Record struct {
	TableName string
	Values    []interface{}
}

func (r Record) GetTableDesc() (*TableDesc, error) {
	for _, t := range tables {
		if t.Name == r.TableName {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unregiester table %s", r.TableName)
}

var tables []TableDesc

// NewStore creates new store implementation.
func NewStore() Store {
	return &store{}
}

func (s *store) RegisterTable(tableDesc TableDesc) error {
	tables = append(tables, tableDesc)
	// TODO handle dumplicate register
	return nil
}

// Save
// 1. the order of values must be same with the order of table desc
func (s *store) Insert(record Record) (savedRecord *Record, err error) {
	tableDesc, err := record.GetTableDesc()
	if err != nil {
		return nil, err
	}

	primaryKey, ptype, columnIndex := tableDesc.GetPrimaryKey()
	if primaryKey != "" {
		if err := s.checkDuplicatedRecord(tableDesc, primaryKey, ptype, record.Values[columnIndex]); err != nil {
			return nil, err
		}
	}

	recordBytes, err := getRecordBytes(record)
	if err != nil {
		return nil, err
	}

	tableFile, err := tableDesc.GetTableFile()
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(tableFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, err = f.Write(recordBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *store) checkDuplicatedRecord(tableDesc *TableDesc, primaryKey string, primaryKeyType ColumnTypes, primaryValue interface{}) error {
	tableFile, err := tableDesc.GetTableFile()
	if err != nil {
		return err
	}

	offsetOfRecord, err := tableDesc.OffsetOfColumn(primaryKey)
	if err != nil {
		return err
	}

	totalBytes := tableDesc.GetTotalBytes()

	f, err := os.OpenFile(tableFile, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		return err
	}

	defer f.Close()

	bs := make([]byte, sizeOf(primaryKeyType))
	offset := int64(offsetOfRecord)

	for {
		_, err := f.ReadAt(bs, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch primaryKeyType {
		case ColumnTypeString:
			var str string
			index := bytes.IndexByte(bs, byte(0))
			if index != -1 {
				str = string(bs[:index])
			} else {
				str = string(bs)
			}
			if primaryValue.(string) == str {
				return ErrDuplicatedRecord
			}
		case ColumnTypeInteger:
			buf := bytes.NewBuffer(bs)
			var x int64
			binary.Read(buf, binary.BigEndian, &x)
			if primaryValue.(int64) == x {
				return ErrDuplicatedRecord
			}
		case ColumnTypeBool:
			return ErrBoolNotSupportPrimaryKey
		}

		offset += int64(totalBytes)
	}

	return nil
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
