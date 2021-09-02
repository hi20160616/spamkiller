package configs

import (
	"os"
	"strings"
)

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
