package main

import (
	"log"
	"os"
	"path/filepath"
)

var LogWriter *os.File

func init() {
	LogWriter, err := os.OpenFile(filepath.Join("./", "log.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Println(err)
	}
	log.SetOutput(LogWriter)
}
