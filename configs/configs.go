package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ProjectName = "spamkiller"
	V           = &configuration{}
	LogWriter   = &os.File{}
)

type configuration struct {
	RootPath    string
	Debug       bool
	Folder      string
	MailSet     string // to copy treated emls
	DropDaysAgo int    `json:"DropDaysAgo"`
	Drop        time.Time
	Filter      struct {
		Spams   []string
		Focuses []string
	} `json:"filter"`
}

func setRootPath() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	V.RootPath = root
	fmt.Println(V.RootPath)

	if strings.Contains(os.Args[0], ".test") {
		return rootPath4Test()
	}
	return nil
}

func load() error {
	cf := filepath.Join(V.RootPath, "configs/configs.json")
	f, err := os.ReadFile(cf)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(f, V); err != nil {
		return err
	}

	// Drop
	V.Drop = time.Now().AddDate(0, 0, -V.DropDaysAgo)

	// load focuses.json
	fp := filepath.Join(V.RootPath, "configs/focuses.json")
	fJson, err := os.ReadFile(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("warning: no focuses.json")
		} else {
			return err
		}
	}
	if err = json.Unmarshal(fJson, &V.Filter.Focuses); err != nil {
		return err
	}

	// load spams.json
	sp := filepath.Join(V.RootPath, "configs/spams.json")
	sJson, err := os.ReadFile(sp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("warning: no spams.json")
		} else {
			return err
		}
	}
	if err = json.Unmarshal(sJson, &V.Filter.Spams); err != nil {
		return err
	}

	// write to configs.json
	data, err := json.MarshalIndent(V, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cf, data, 0755)
}

func init() {
	LogWriter, err := os.OpenFile(filepath.Join("./", "log.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Println(err)
	}
	log.SetOutput(LogWriter)
	if err := setRootPath(); err != nil {
		log.Printf("configs init error: %v\n", err)
	}
	if err := load(); err != nil {
		log.Printf("configs load error: %v\n", err)
	}
}

func rootPath4Test() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	ps := strings.Split(root, ProjectName)
	n := 0
	if len(ps) > 1 {
		n = strings.Count(ps[1], string(os.PathSeparator))
	}
	for i := 0; i < n; i++ {
		V.RootPath = filepath.Join("../", "./")
	}
	return nil
}
