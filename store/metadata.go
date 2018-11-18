package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ColumnTypes int

const (
	_ ColumnTypes = iota
	ColumnTypeString
	ColumnTypeBool
	ColumnTypeInteger
	ColumnTypeByte
)

type Column struct {
	Name       string
	Type       ColumnTypes
	PrimaryKey bool
}

type TableDesc struct {
	Name       string
	Columns    []Column
	MaxInnerID int64
}

var tables []TableDesc

// clearTableDescs clears table descs, only used in test
func clearTableDescs() {
	tables = make([]TableDesc, 0)
}

type Record struct {
	TableName string
	Values    []interface{}
}

type UpdateSetValueFn func(record Record) (interface{}, error)

type SetItem struct {
	Name  string
	Value UpdateSetValueFn
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

func (t TableDesc) IndexOfColumn(columnName string) (int, error) {
	for i, c := range t.Columns {
		if c.Name == columnName {
			return i, nil
		}
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

func (t TableDesc) GetTableFile() (string, error) {
	workingPath, err := getWorkingPath()
	if err != nil {
		return "", err
	}

	tableFile := filepath.Join(workingPath, t.Name)
	return tableFile, nil
}

func (r Record) GetTableDesc() (*TableDesc, error) {
	return GetTableDescFromTableName(r.TableName)
}

func GetTableDescFromTableName(tableName string) (*TableDesc, error) {
	for _, t := range tables {
		if t.Name == tableName {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unregiester table %s", tableName)
}

func sizeOf(t ColumnTypes) int {
	switch t {
	case ColumnTypeString:
		return 255
	case ColumnTypeInteger:
		return 8
	case ColumnTypeBool:
		return 1
	case ColumnTypeByte:
		return 1
	}

	panic("unsupport type")
}

func SaveMetadata() error {
	descBytes, err := json.Marshal(tables)
	if err != nil {
		panic(err)
	}
	metaDatafilePath, err := getMetadataFilePath()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(metaDatafilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteAt(descBytes, 0)
	if err != nil {
		return err
	}

	return nil
}

func LoadMetadata() error {
	metaDatafilePath, err := getMetadataFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(metaDatafilePath); os.IsNotExist(err) {
		return nil
	}

	metadata, err := ioutil.ReadFile(metaDatafilePath)
	if err != nil {
		return err
	}

	tableDescs := []TableDesc{}
	err = json.Unmarshal(metadata, &tableDescs)
	if err != nil {
		return err
	}

	tables = tableDescs

	// TODO make sure MaxInnerID in disk is right
	return nil
}

func GetMetadataOf(tableName string) (*TableDesc, error) {
	for _, t := range tables {
		if t.Name == tableName {
			outDesc := t
			outDesc.Columns = t.Columns[2:]
			return &outDesc, nil
		}
	}

	return nil, errors.New("%s not found")

}

func getMetadataFilePath() (path string, err error) {
	workingPath, err := getWorkingPath()
	if err != nil {
		return "", err
	}

	metadataFilePath := filepath.Join(workingPath, "____metadata____")
	return metadataFilePath, nil
}
