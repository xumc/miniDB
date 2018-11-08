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

	store.NewStore(logger)

	// find students who didn't pass the exam, age < 18 and female
	qt := store.QueryTree{
		MatchAll: true,
		Negative: false,
		Left: &store.QueryTree{
			MatchAll: true,
			Negative: false,

			Left: &store.QueryTree{
				Item: &store.QueryItem{Key: "pass", Operator: store.MatcherEqual{}, Value: false},
			},
			Right: &store.QueryTree{
				Item: &store.QueryItem{Key: "age", Operator: store.MatcherLessThan{}, Value: int64(18)},
			},
		},
		Right: &store.QueryTree{
			Item: &store.QueryItem{Key: "sex", Operator: store.MatcherEqual{}, Value: "FEMALE"},
		},
	}
	fmt.Println(qt.PrettyPrint())

	fmt.Println("finish")
}
