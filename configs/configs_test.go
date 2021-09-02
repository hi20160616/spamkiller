package configs

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg := &Config{ProjectName: "spamkiller"}
	if err := setRootPath(cfg).load(); err != nil {
		t.Error(err)
	}
	fmt.Println(cfg.Drop)
	for _, e := range cfg.Filter.Spams {
		fmt.Println(e)
	}
	for _, e := range cfg.Filter.Focuses {
		fmt.Println(e)
	}
}

func TestRootPath4Test(t *testing.T) {
	cfg := rootPath4Test(&Config{})
	if cfg.Err != nil {
		t.Error(cfg.Err)
	}
	fmt.Println(cfg.RootPath)
}
