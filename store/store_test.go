package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestCRUD(t *testing.T) {
	clearTestData()

	fileName := "miniDB.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	logger := log.New(logFile, "[Debug]", log.LstdFlags)

	s := NewStore(logger)

	fmt.Println("----------------regiester------------------")

	err = s.RegisterTable(TableDesc{
		Name: "student",
		Columns: []Column{
			Column{Name: "id", Type: ColumnTypeInteger, PrimaryKey: true},
			Column{Name: "name", Type: ColumnTypeString},
			Column{Name: "pass", Type: ColumnTypeBool},
		},
	})

	if err != nil {
		t.Fatalf("err should be nil, but %s", err)
	}

	fmt.Println("----------------insert------------------")
	r1 := Record{
		TableName: "student",
		Values:    []interface{}{int64(1), "Jack 1", true},
	}
	r2 := Record{
		TableName: "student",
		Values:    []interface{}{int64(2), "Jack 2", false},
	}
	insertRecord(s, r1)
	insertRecord(s, r2)
	assertRecords(s, t, [][]interface{}{
		[]interface{}{int64(1), "Jack 1", true},
		[]interface{}{int64(2), "Jack 2", false},
	})

	fmt.Println("----------------update------------------")

	updateFn := func(r Record) (interface{}, error) {
		return "prefix_" + r.Values[1].(string), nil
	}

	updateRecords(
		s,
		&QueryTree{
			Item: &QueryItem{Key: "id", Operator: MatcherEqual{}, Value: int64(2)},
		},
		[]UpdateSetItem{
			UpdateSetItem{Name: "name", Value: updateFn},
		},
	)
	assertRecords(s, t, [][]interface{}{
		[]interface{}{int64(1), "Jack 1", true},
		[]interface{}{int64(2), "prefix_Jack 2", false},
	})

	fmt.Println("----------------delete------------------")
	deleteRecords(s, &QueryTree{
		Item: &QueryItem{Key: "pass", Operator: MatcherEqual{}, Value: true},
	})
	assertRecords(s, t, [][]interface{}{
		[]interface{}{int64(2), "prefix_Jack 2", false},
	})
}

func insertRecord(s Store, record Record) {
	affectedRows, err := s.Insert(record.TableName, record)
	if err != nil {
		if err == ErrDuplicatedRecord {
			fmt.Println("dumplicated record found in db")
			return
		}
		if err == ErrBoolNotSupportPrimaryKey {
			fmt.Println("bool type doesn't support primary key setting")
		}

		panic(err)
	}

	fmt.Println("insert : ", affectedRows)
}

func assertRecords(s Store, t *testing.T, expected [][]interface{}) {
	rs, err := s.Select("student", nil)
	if err != nil {
		panic(err)
	}

	for i, r := range rs {
		for j := range r.Values {
			if rs[i].Values[j] != expected[i][j] {
				t.Fatalf("value should be %v, but got %v", expected[i][j], rs[i].Values[j])
			}
		}
	}
	return
}

func deleteRecords(s Store, qt *QueryTree) {
	affectedRows, err := s.Delete("student", qt)
	if err != nil {
		panic(err)
	}

	fmt.Println("deleted: ", affectedRows)
}

func updateRecords(s Store, qt *QueryTree, uis []UpdateSetItem) {
	affectedRows, err := s.Update("student", qt, uis)
	if err != nil {
		panic(err)
	}

	fmt.Println("updated: ", affectedRows)
}

func clearTestData() {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	sfile := strings.Split(file, "/")
	datafile := strings.Join(append(sfile[:(len(sfile)-2)], "student"), "/")

	os.Remove(datafile)
}
