package main

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	w, err := os.OpenFile(filepath.Join("./", "log.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Println(err)
	}
	log.SetOutput(w)
	defer w.Close()
}
