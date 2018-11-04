package main

import (
	"fmt"

	"github.com/xumc/miniDB/store"
)

func main() {
	s := store.NewStore()

	s.RegisterTable(store.TableDesc{
		Name: "student",
		Columns: []store.Column{
			store.Column{Name: "id", Type: store.ColumnTypeInteger},
			store.Column{Name: "name", Type: store.ColumnTypeString, PrimaryKey: true},
			store.Column{Name: "pass", Type: store.ColumnTypeBool},
		},
	})

	record := store.Record{
		TableName: "student",
		Values:    []interface{}{int64(5), "Jack 3", true},
	}

	savedRecord, err := s.Insert(record)
	if err != nil {
		if err == store.ErrDuplicatedRecord {
			fmt.Println("dumplicated record found in db")
			return
		}
		if err == store.ErrBoolNotSupportPrimaryKey {
			fmt.Println("bool type doesn't support primary key setting")
		}

		panic(err)
	}

	fmt.Println(savedRecord)
	fmt.Println("finish")
}
