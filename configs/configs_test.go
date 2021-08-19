package configs

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	if err := load(); err != nil {
		t.Error(err)
	}
	for _, e := range V.Filter.Spams {
		fmt.Println(e)
	}
	for _, e := range V.Filter.Focuses {
		fmt.Println(e)
	}
}

func TestRootPath4Test(t *testing.T) {
	if err := rootPath4Test(); err != nil {
		t.Error(err)
	}
	fmt.Println(V.RootPath)
}
