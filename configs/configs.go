//go:build linux||darwin

package configs

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	ProjectName string
	RootPath    string
	Raw         []byte
	Debug       bool
	Verbose     bool // if true, prompt enter to exit.
	LogName     string
	MailTypes   []string
	Folder      string
	Result      string // to copy treated emls
	DropDaysAgo int    `json:"DropDaysAgo"`
	Drop        time.Time
	Filter      struct {
		Spams   []string
		Focuses []string
	} `json:"filter"`
	Err error
}

func NewConfig(projectName string) *Config {
	return setRootPath(&Config{ProjectName: projectName}).load()
}

func rootPath4Test(cfg *Config) *Config {
	cfg.RootPath, cfg.Err = os.Getwd()
	if cfg.Err != nil {
		return cfg
	}
	ps := strings.Split(cfg.RootPath, cfg.ProjectName)
	n := 0
	if len(ps) > 1 {
		n = strings.Count(ps[1], string(os.PathSeparator))
	}
	for i := 0; i < n; i++ {
		cfg.RootPath = filepath.Join("../", "./")
	}
	return cfg
}

func (c *Config) load() *Config {
	if c.Err != nil {
		return c
	}
	cfgFile := filepath.Join(c.RootPath, "configs", "configs.json")
	c.Raw, c.Err = os.ReadFile(cfgFile)
	if c.Err != nil {
		if errors.Is(c.Err, os.ErrNotExist) {
			c.Err = errors.WithMessage(c.Err, "ReadFile error: no configs.json")
		}
		return c
	}
	cfgTemp := &Config{}
	if c.Err = json.Unmarshal(c.Raw, cfgTemp); c.Err != nil {
		c.Err = errors.WithMessage(c.Err, "Unmarshal configs.json error")
		return c
	}
	c.Debug = cfgTemp.Debug
	c.Verbose = cfgTemp.Verbose
	c.LogName = cfgTemp.LogName
	c.ProjectName = cfgTemp.ProjectName
	c.Folder = cfgTemp.Folder
	c.Result = cfgTemp.Result
	c.MailTypes = cfgTemp.MailTypes
	c.DropDaysAgo = cfgTemp.DropDaysAgo

	// Drop
	c.Drop = time.Now().AddDate(0, 0, -c.DropDaysAgo)

	// load *.json
	loadJson := func(filename string) ([]string, error) {
		fp := filepath.Join(c.RootPath, "configs", filename)
		fJson, err := os.ReadFile(fp)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Println("warning: no ", filename)
			} else {
				return nil, err
			}
		}
		keywords := []string{}
		if err = json.Unmarshal(fJson, &keywords); err != nil {
			return nil, errors.WithMessagef(err, "Unmarshal %s error", filename)
		}
		return keywords, nil
	}

	// load focuses.json
	c.Filter.Focuses, c.Err = loadJson("focuses.json")
	// load spams.json
	c.Filter.Spams, c.Err = loadJson("spams.json")
	return c
}

func setRootPath(cfg *Config) *Config {
	cfg.RootPath, cfg.Err = os.Getwd()
	if cfg.Err != nil {
		return cfg
	}
	if strings.Contains(os.Args[0], ".test") {
		return rootPath4Test(cfg)
	}
	return cfg
}
