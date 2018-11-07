package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xumc/miniDB/store"
)

func main() {
	fileName := "miniDB.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	logger := log.New(logFile, "[Debug]", log.LstdFlags)

	s := store.NewStore(logger)

	s.RegisterTable(store.TableDesc{
		Name: "student",
		Columns: []store.Column{
			store.Column{Name: "id", Type: store.ColumnTypeInteger, PrimaryKey: true},
			store.Column{Name: "name", Type: store.ColumnTypeString},
			store.Column{Name: "pass", Type: store.ColumnTypeBool},
		},
	})

	fmt.Println("----------------insert------------------")
	r1 := store.Record{
		TableName: "student",
		Values:    []interface{}{int64(1), "Jack 1", true},
	}
	r2 := store.Record{
		TableName: "student",
		Values:    []interface{}{int64(2), "Jack 2", false},
	}
	insertRecord(s, r1)
	insertRecord(s, r2)
	selectRecords(s, []store.QueryItem{})

	fmt.Println("----------------update------------------")

	updateFn := func(r store.Record) (interface{}, error) {
		return "prefix_" + r.Values[1].(string), nil
	}

	updateRecords(
		s,
		[]store.QueryItem{
			store.QueryItem{Key: "id", Operator: store.QueryOperatorEqual, Value: int64(2)},
		},
		[]store.UpdateSetItem{
			store.UpdateSetItem{Name: "name", Value: updateFn},
		},
	)
	selectRecords(s, []store.QueryItem{})

	fmt.Println("----------------delete------------------")
	deleteRecords(s, []store.QueryItem{
		store.QueryItem{Key: "pass", Operator: store.QueryOperatorEqual, Value: true},
	})
	selectRecords(s, []store.QueryItem{})

	fmt.Println("finish")
}

func updateRecords(s store.Store, qs []store.QueryItem, uis []store.UpdateSetItem) {
	affectedRows, err := s.Update("student", qs, uis)
	if err != nil {
		panic(err)
	}

	fmt.Println("updated: ", affectedRows)
}

func deleteRecords(s store.Store, qs []store.QueryItem) {
	affectedRows, err := s.Delete("student", qs)
	if err != nil {
		panic(err)
	}

	fmt.Println("deleted: ", affectedRows)
}

func selectRecords(s store.Store, qs []store.QueryItem) {
	rs, err := s.Select("student", qs)
	if err != nil {
		panic(err)
	}

	for _, r := range rs {
		fmt.Println(r.Values)
	}

	return
}

func insertRecord(s store.Store, record store.Record) {
	affectedRows, err := s.Insert(record.TableName, record)
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

	fmt.Println("insert : ", affectedRows)
}
