package configs

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	ProjectName = "spamkiller"
	V           = &configuration{}
)

type configuration struct {
	RootPath string
	Debug    bool
	Folder   string
	Filter   struct {
		Spams   []string
		Focuses []string
	} `json:"filter"`
}

func setRootPath() error {
	V.RootPath = "./"
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
	// load focuses.json
	fJson, err := os.ReadFile(filepath.Join(V.RootPath, "configs/focuses.json"))
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
	sJson, err := os.ReadFile(filepath.Join(V.RootPath, "configs/spams.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("warning: no focuses.json")
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
	if err := setRootPath(); err != nil {
		log.Printf("configs init error: %v", err)
	}
	if err := load(); err != nil {
		log.Printf("configs load error: %v", err)
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
		V.RootPath = filepath.Join("../", V.RootPath)
	}
	return nil
}
