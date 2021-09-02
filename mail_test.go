package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hi20160616/spamkiller/configs"
)

var testCfg = configs.NewConfig("spamkiller")
var testLog = func() *log.Logger {
	log, err := NewLog(testCfg)
	if err != nil {
		fmt.Println(err)
	}
	return log
}()

var m *Mail = func() *Mail {
	testCfg.Result = "./test"
	a := NewMail(
		context.Background(), testCfg, testLog, "./test/test_common.eml")
	if a.err != nil {
		fmt.Println(a.err)
	}
	return a
}()

func TestTreat(t *testing.T) {
	if err := treat(
		context.Background(), testCfg, testLog, "./test"); err != nil {
		t.Error(err)
	}
}

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

	defer configs.LogWriter.Close()
}

func TestAnalysis(t *testing.T) {
	// condition 0
	m := m.analysis()
	if m.tag != 0 {
		t.Errorf("want: 1, got: %d", m.tag)
	}

	// condition 1
	m.cfg.Filter.Focuses = append(m.cfg.Filter.Focuses, "hi2020")
	m = m.analysis()
	if m.tag != 1 {
		t.Errorf("want: 1, got: %d", m.tag)
	}

	// condition 3
	m.cfg.Filter.Spams = append(m.cfg.Filter.Spams, "outlook.com")
	m = m.analysis()
	if m.tag != 3 {
		t.Errorf("want: 3, got: %d", m.tag)
	}

	// condition 2
	m.cfg.Filter.Focuses = m.cfg.Filter.Focuses[:len(m.cfg.Filter.Focuses)-1]
	m = m.analysis()
	if m.tag != 2 {
		t.Errorf("want: 2, got: %d", m.tag)
	}

	defer configs.LogWriter.Close()
}

var ms *Mails = NewMails(context.Background(), testCfg, testLog, "./test")

func TestWalkSrc(t *testing.T) {
	ms = ms.walkSrc(context.Background())
	if ms.err != nil {
		t.Error(ms.err)
	}
}
