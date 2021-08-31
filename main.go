package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/hi20160616/spamkiller/configs"
)

func main() {
	defer configs.LogWriter.Close()
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	configs.V.Folder = os.Args[1]
	if configs.V.Result == "./" {
		configs.V.Result = os.Args[1]
	}
	if err := treat(configs.V.Folder); err != nil {
		fmt.Println(err)
		log.Println(err)
		fmt.Println("Sth error. Press Enter to quit!")
		// bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Scanln()
	}
	if configs.V.Verbose {
		fmt.Println("Done. Press Enter to quit!")
		// bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Scanln()
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
		if filepath.Ext(ep) != ".eml" {
			continue
		}
		// new mail
		m := NewMail(ep)
		if m.err != nil {
			log.Println(err)
		}
		// analysis and deliver, log out the error
		if err := m.analysis().deliver(); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func getEmlPathes(scanPath string) ([]string, error) {
	subDirToSkip := "skip"
	emls := make([]string, 0)
	if err := filepath.WalkDir(scanPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if d.IsDir() && d.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", d.Name())
			return filepath.SkipDir
		}
		if d.Type().IsRegular() {
			emls = append(emls, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return emls, nil
}
