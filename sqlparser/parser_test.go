package sqlparser

import (
	"log"
	"testing"

	"github.com/kr/pretty"

	"github.com/stretchr/testify/assert"
)

var (
	einsert    = "INSERT"
	eupdate    = "UPDATE"
	eselect    = "SELECT"
	edelete    = "DELETE"
	etableName = "student"
	eid        = "id"
	ename      = "name"
	eage       = "age"
	exumc      = "xumc"
	emxu       = "mxu"
	e1         = int64(1)
	e30        = int64(30)
	etrue      = true
	eequal     = Operator("=")
	epass      = "pass"
	ebtrue     = Boolean(true)
)

func TestInsert(t *testing.T) {
	parser := NewParser(&log.Logger{})

	sql, err := parser.Parse("INSERT INTO student(id, name, age, pass) VALUES(1, \"xumc\", 30, true);")
	assert.Empty(t, err)
	isql := sql.(*InsertSQL)

	expected := InsertSQL{
		Type:      &einsert,
		TableName: &etableName,
		Fields: []*InsertField{
			&InsertField{
				Name: &eid,
			},
			&InsertField{
				Name: &ename,
			},
			&InsertField{
				Name: &eage,
			},
			&InsertField{
				Name: &epass,
			},
		},
		Values: []*InsertValue{
			&InsertValue{
				String: (*string)(nil),
				Number: &e1,
			},
			&InsertValue{
				String: &exumc,
				Number: (*int64)(nil),
			},
			&InsertValue{
				String: (*string)(nil),
				Number: &e30,
			},
			&InsertValue{
				String:  (*string)(nil),
				Number:  (*int64)(nil),
				Boolean: &ebtrue,
			},
		},
	}
	diffs := pretty.Diff(expected, *isql)
	assert.Equal(t, 0, len(diffs))
}

func TestUpdate(t *testing.T) {
	parser := NewParser(&log.Logger{})

	sql, err := parser.Parse("UPDATE student SET name=\"mxu\", age=30 WHERE id=1;")
	assert.Empty(t, err)
	usql := sql.(*UpdateSQL)

	expected := UpdateSQL{
		Type:      &eupdate,
		TableName: &etableName,
		SetItems: []*SetItem{
			&SetItem{
				Key: &ename,
				Value: &SetValue{
					String: &emxu,
					Number: (*int64)(nil),
				},
			},
			&SetItem{
				Key: &eage,
				Value: &SetValue{
					String: (*string)(nil),
					Number: &e30,
				},
			},
		},
		Where: &QueryTree{
			Negative:  false,
			LeftTree:  (*QueryTree)(nil),
			MatchAll:  false,
			RightTree: (*QueryTree)(nil),
			Item: &QueryItem{
				Key:      &eid,
				Operator: &eequal,
				Value: &QueryValue{
					String: (*string)(nil),
					Number: &e1,
				},
			},
		},
	}

	diffs := pretty.Diff(expected, *usql)
	assert.Equal(t, 0, len(diffs))
}

func TestSelect(t *testing.T) {
	parser := NewParser(&log.Logger{})

	sql, err := parser.Parse("SELECT name FROM student WHERE id=1;")
	assert.Empty(t, err)
	ssql := sql.(*SelectSQL)

	expected := SelectSQL{
		Type: &eselect,
		Fields: []*FieldValue{
			&FieldValue{
				String:    &ename,
				AllFields: false,
			},
		},
		TableName: &etableName,
		Where: &QueryTree{
			Negative:  false,
			LeftTree:  (*QueryTree)(nil),
			MatchAll:  false,
			RightTree: (*QueryTree)(nil),
			Item: &QueryItem{
				Key:      &eid,
				Operator: &eequal,
				Value: &QueryValue{
					String: (*string)(nil),
					Number: &e1,
				},
			},
		},
	}

	diffs := pretty.Diff(expected, *ssql)
	assert.Equal(t, 0, len(diffs))
}

func TestDelete(t *testing.T) {
	parser := NewParser(&log.Logger{})

	sql, err := parser.Parse("DELETE FROM student WHERE id=1;")
	assert.Empty(t, err)
	dsql := sql.(*DeleteSQL)

	expected := DeleteSQL{
		Type:      &edelete,
		TableName: &etableName,
		Where: &QueryTree{
			Negative:  false,
			LeftTree:  (*QueryTree)(nil),
			MatchAll:  false,
			RightTree: (*QueryTree)(nil),
			Item: &QueryItem{
				Key:      &eid,
				Operator: &eequal,
				Value: &QueryValue{
					String: (*string)(nil),
					Number: &e1,
				},
			},
		},
	}
	diffs := pretty.Diff(expected, *dsql)
	assert.Equal(t, 0, len(diffs))
}
