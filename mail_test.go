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
		fmt.Println("testLog error:", err)
	}
	return log
}()

var m *Mail = func() *Mail {
	testCfg.Result = "./test"
	testCfg.Folder = "./test"
	a := NewMail(
		context.Background(), testCfg, testLog, "./test/test_common.eml")
	if a.err != nil {
		fmt.Println(a.err)
	}
	return a
}()

func TestAnalysis(t *testing.T) {
	// condition 2 FOCUSED
	m.cfg.Filter.Focuses = append(m.cfg.Filter.Focuses, "hi2020")
	m = m.analysis()
	if m.tag != 2 {
		t.Errorf("want: 2, got: %d", m.tag)
	}

	// condition 1 COMMON
	m.cfg.Filter.Focuses = m.cfg.Filter.Focuses[:len(m.cfg.Filter.Focuses)-1]
	m = m.analysis()
	if m.tag != 1 {
		t.Errorf("want: 1, got: %d", m.tag)
	}

	// condition 3 SPAM
	m.cfg.Filter.Spams = append(m.cfg.Filter.Spams, "outlook.com")
	m = m.analysis()
	if m.tag != 3 {
		t.Errorf("want: 3, got: %d", m.tag)
	}

	// condition 0 ERR
	m := NewMail(context.Background(), testCfg, testLog, "./test/testErrMail.eml").analysis()
	if m.tag != 0 {
		t.Errorf("want: 0, got: %d", m.tag)
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
	m.tag = 4
	if err := m.deliver(); err != nil {
		t.Error(err)
	}
	m.tag = 5
	if err := m.deliver(); err != nil {
		t.Error(err)
	}
}

func TestTreat(t *testing.T) {
	if err := treat(
		context.Background(), testCfg, testLog, "./test/testErrMail.eml"); err != nil {
		t.Error(err)
	}
	if err := treat(
		context.Background(), testCfg, testLog, "./test/test_common.eml"); err != nil {
		t.Error(err)
	}
	if err := treat(
		context.Background(), testCfg, testLog, "./test/just a test.eml"); err != nil {
		t.Error(err)
	}
}

var ms *Mails = NewMails(context.Background(), testCfg, testLog, "./test")

func TestWalkSrc(t *testing.T) {
	// FOCUSED
	ms.cfg.Filter.Focuses = append(ms.cfg.Filter.Focuses, "hi2020")
	ms = ms.walkSrc(context.Background())
	if ms.err != nil {
		t.Error(ms.err)
	}
	// undo
	ms.cfg.Filter.Focuses = ms.cfg.Filter.Focuses[:len(ms.cfg.Filter.Focuses)-1]

	// SPAM
	ms.cfg.Filter.Spams = append(ms.cfg.Filter.Spams, "outlook.com")
	ms = ms.walkSrc(context.Background())
	if ms.err != nil {
		t.Error(ms.err)
	}
	// undo
	ms.cfg.Filter.Spams = ms.cfg.Filter.Focuses[:len(ms.cfg.Filter.Spams)-1]

	// CONFUSION
	ms.cfg.Filter.Focuses = append(ms.cfg.Filter.Focuses, "hi2020")
	ms.cfg.Filter.Spams = append(ms.cfg.Filter.Spams, "outlook.com")
	ms = ms.walkSrc(context.Background())
	if ms.err != nil {
		t.Error(ms.err)
	}
}
