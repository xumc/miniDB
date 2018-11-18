package main

import (
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

	err = store.LoadMetadata()
	if err != nil {
		panic(err)
	}

	store.NewStore(logger)
}
