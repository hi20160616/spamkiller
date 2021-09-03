package configs

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func setRootPath(cfg *Config) *Config {
	cfg.RootPath, cfg.Err = readKey()
	if cfg.Err != nil {
		return cfg
	}
	cfg.RootPath = filepath.Dir(cfg.RootPath)
	if strings.Contains(os.Args[0], ".test") {
		return rootPath4Test(cfg)
	}
	return cfg
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
	s = strings.Split(s, " ")[0]
	return strings.ReplaceAll(s, "\"", ""), nil
}
