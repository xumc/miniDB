package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xumc/miniDB/connection"
	"github.com/xumc/miniDB/sqlparser"
	"github.com/xumc/miniDB/transaction"

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

	t := transaction.NewTransaction(logger, s)

	p := sqlparser.NewParser(logger, t)

	c := connection.NewConnection(logger, p)

	c.Run()

	// var g group.Group

	// storge component runs
	// {
	// 	g.Add(
	// 		func() error {
	// 			s.Run()
	// 			return nil
	// 		},
	// 		func(err error) {
	// 		},
	// 	)
	// }

	// transaction component runs
	// {
	// 	g.Add(
	// 		func() error {
	// 			t.Run()
	// 			return nil
	// 		},
	// 		func(err error) {
	// 		},
	// 	)
	// }

	// parser component runs
	// {
	// 	g.Add(
	// 		func() error {
	// 			p.Run()
	// 			return nil
	// 		},
	// 		func(err error) {
	// 		},
	// 	)
	// }

	// connection component runs
	// {
	// 	g.Add(
	// 		func() error {
	// 			c.Run()
	// 			return nil
	// 		},
	// 		func(err error) {
	// 		},
	// 	)
	// }

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
