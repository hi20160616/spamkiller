package main

import "testing"

func TestTreat(t *testing.T) {
	if err := treat("./test"); err != nil {
		t.Error(err)
	}
}
