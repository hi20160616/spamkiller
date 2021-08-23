package configs

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/sys/windows/registry"
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
	if runtime.GOOS == "Linux" {
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		V.RootPath = root
	}
	if runtime.GOOS == "windows" {
		s, err := readKey()
		if err != nil {
			return err
		}
		V.RootPath = filepath.Dir(strings.ReplaceAll(s, "\"", ""))
	}
	if strings.Contains(os.Args[0], ".test") {
		return rootPath4Test()
	}
	return nil
}

func readKey() (string, error) {
	k, err := registry.OpenKey(registry.CLASSES_ROOT, `Directory\shell\SpamKiller\command`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("")
	if err != nil {
		return "", err
	}
	return strings.Split(s, " ")[0], nil
}

func load() error {
	cf := filepath.Join(V.RootPath, "configs", "configs.json")
	f, err := os.ReadFile(cf)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.WithMessage(err, "ReadFile error: no configs.json")
		} else {
			return err
		}
	}
	if err = json.Unmarshal(f, V); err != nil {
		return errors.WithMessage(err, "Unmarshal configs.json error")
	}

	// Drop
	V.Drop = time.Now().AddDate(0, 0, -V.DropDaysAgo)

	// load focuses.json
	fp := filepath.Join(V.RootPath, "configs", "focuses.json")
	fJson, err := os.ReadFile(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("warning: no focuses.json")
		} else {
			return err
		}
	}
	if err = json.Unmarshal(fJson, &V.Filter.Focuses); err != nil {
		return errors.WithMessage(err, "Unmarshal focuses.json error")
	}

	// load spams.json
	sp := filepath.Join(V.RootPath, "configs", "spams.json")
	sJson, err := os.ReadFile(sp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("warning: no spams.json")
		} else {
			return errors.WithMessage(err, "Unmarshal spams.json error")
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
		log.Printf("configs init error: %v\n", err)
	}
	if err := load(); err != nil {
		log.Printf("configs load error: %v\n", err)
	}
	LogWriter, err := os.OpenFile(filepath.Join(V.RootPath, "log.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Println(err)
	}
	log.SetOutput(LogWriter)
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
