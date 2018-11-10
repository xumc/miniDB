package store

import (
	"fmt"
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
	Name    string
	Columns []Column
}

type Record struct {
	TableName string
	Values    []interface{}
}

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

type UpdateSetValueFn func(record Record) (interface{}, error)

type UpdateSetItem struct {
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
