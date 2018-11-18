package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xumc/miniDB/sqlparser"

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

	err = store.LoadMetadata()
	if err != nil {
		panic(err)
	}

	s := store.NewStore(logger)

	p := sqlparser.NewParser(logger)

	// sql, err := p.Parse("INSERT INTO student(id, name, pass) VALUES(1, \"xumcxumc\", true);")
	// sql, err := p.Parse("UPDATE student SET name=\"xumcxumcxumc\" WHERE name=\"xumcxumc\";")
	sql, err := p.Parse("SELECT * FROM student WHERE name=\"xumcxumc\";")
	// sql, err := p.Parse("DELETE FROM student WHERE name=\"xumcxumc\";")
	if err != nil {
		panic(err)
	}

	switch sqlStruct := sql.(type) {
	case *sqlparser.InsertSQL:
		tableDesc, err := store.GetMetadataOf(*sqlStruct.TableName)
		if err != nil {
			panic(err)
		}

		record := p.TransformInsert(sqlStruct, tableDesc)

		_, err = s.Insert(record.TableName, record)
		if err != nil {
			panic(err)
		}
	case *sqlparser.UpdateSQL:
		tableDesc, err := store.GetMetadataOf(*sqlStruct.TableName)
		if err != nil {
			panic(err)
		}

		qt, setItems := p.TransformUpdate(sqlStruct, tableDesc)

		_, err = s.Update(*sqlStruct.TableName, qt, setItems)
		if err != nil {
			panic(err)
		}
	case *sqlparser.SelectSQL:
		qt := p.TransformSelect(sqlStruct)

		rs, err := s.Select(*sqlStruct.TableName, qt)
		if err != nil {
			panic(err)
		}

		fmt.Println(rs)
	case *sqlparser.DeleteSQL:
		qt := p.TransformDelete(sqlStruct)

		rs, err := s.Delete(*sqlStruct.TableName, qt)
		if err != nil {
			panic(err)
		}

		fmt.Println(rs)
	default:
	}

	printRecords(s)

	fmt.Println("done")
}

func printRecords(s store.Store) {
	rs, err := s.Select("student", nil)
	if err != nil {
		panic(err)
	}

	for i, r := range rs {
		for j := range r.Values {
			fmt.Printf("%v 	", rs[i].Values[j])
		}
		fmt.Println("")
	}
}
