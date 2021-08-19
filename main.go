package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/hi20160616/spamkiller/configs"
)

func main() {
	if len(os.Args) == 2 {
		configs.V.Folder = os.Args[1]
	}
	if err := treat(configs.V.Folder); err != nil {
		fmt.Println(err)
		log.Println(err)
	} else {
		fmt.Println("Done. Press Enter to quit!")
	}
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if input.Text() == "\n" {
			return
		}
	}
}

func treat(scanPath string) error {
	// 1. get all emls path
	emlPathes, err := getEmlPathes(scanPath)
	if err != nil {
		return err
	}
	// 2. scan emls and compare emls and filter
	for _, ep := range emlPathes {
		// new mail
		m, err := NewMail(ep)
		if err != nil {
			log.Println(err)
		}
		// analysis and deliver, log out the error
		log.Println(m.analysis().deliver())
	}
	return nil
}

func getEmlPathes(scanPath string) ([]string, error) {
	subDirToSkip := "skip"
	emls := make([]string, 0)
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if d.IsDir() && d.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", d.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		emls = append(emls, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return emls, nil
}
