package main

import (
	"fmt"
	"testing"

	"github.com/hi20160616/spamkiller/configs"
)

var m *Mail = func() *Mail {
	configs.V.MailSet = "./test"
	a, err := NewMail("./test/test_common.eml")
	if err != nil {
		fmt.Print(err)
	}
	return a
}()

func TestDeliver(t *testing.T) {
	m.tag = 0
	if err := m.deliver(); err != nil {
		t.Error(err)
	}

	m.tag = 1
	if err := m.deliver(); err != nil {
		t.Error(err)
	}

	m.tag = 2
	if err := m.deliver(); err != nil {
		t.Error(err)
	}

	m.tag = 3
	if err := m.deliver(); err != nil {
		t.Error(err)
	}

	defer LogWriter.Close()
}

func TestAnalysis(t *testing.T) {
	// condition 0
	m := m.analysis()
	if m.tag != 0 {
		t.Errorf("want: 1, got: %d", m.tag)
	}

	// condition 1
	configs.V.Filter.Focuses = append(configs.V.Filter.Focuses, "hi2020")
	m = m.analysis()
	if m.tag != 1 {
		t.Errorf("want: 1, got: %d", m.tag)
	}

	// condition 3
	configs.V.Filter.Spams = append(configs.V.Filter.Spams, "outlook.com")
	m = m.analysis()
	if m.tag != 3 {
		t.Errorf("want: 3, got: %d", m.tag)
	}

	// condition 2
	configs.V.Filter.Focuses = configs.V.Filter.Focuses[:len(configs.V.Filter.Focuses)-1]
	m = m.analysis()
	if m.tag != 2 {
		t.Errorf("want: 2, got: %d", m.tag)
	}

	defer LogWriter.Close()
}
