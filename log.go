package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hi20160616/spamkiller/configs"
)

func NewLog(cfg *configs.Config) (*log.Logger, error) {
	logWriter, err := os.OpenFile(
		filepath.Join(cfg.RootPath, cfg.LogName),
		os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}
	return log.New(logWriter, "[SK] ", log.LstdFlags), nil
}
