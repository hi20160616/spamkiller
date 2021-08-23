package main

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hi20160616/spamkiller/configs"
)

const (
	COMMON int = iota
	FOCUSED
	SPAM
	CONFUSION
	DROP
)

type Mail struct {
	path          string
	raw           []byte
	sraw          string
	msg           *mail.Message
	date          time.Time
	from          *mail.Address
	to            []*mail.Address
	subject, body string
	tag           int
}

func NewMail(mailPath string) (*Mail, error) {
	r, err := os.ReadFile(mailPath)
	if err != nil {
		return nil, err
	}
	msg, err := mail.ReadMessage(bytes.NewReader(r))
	if err != nil {
		return nil, err
	}
	subject := msg.Header.Get("Subject")
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}
	date, err := msg.Header.Date()
	if err != nil {
		return nil, err
	}
	from, err := msg.Header.AddressList("From")
	if err != nil {
		return nil, err
	}
	to, err := msg.Header.AddressList("To")
	if err != nil {
		return nil, err
	}
	sraw := msg.Header.Get("From") + msg.Header.Get("To") +
		subject + " " + string(body)
	return &Mail{
		path:    mailPath,
		raw:     r,
		sraw:    sraw,
		msg:     msg,
		date:    date,
		from:    from[0],
		to:      to,
		subject: subject,
		body:    string(body),
	}, nil
}

func (m *Mail) analysis() *Mail {
	flag := COMMON
	if m.date.Before(configs.V.Drop) {
		m.tag = DROP
		return m
	}
	for _, kw := range configs.V.Filter.Spams {
		if strings.Contains(m.sraw, kw) {
			flag = SPAM
		}
	}
	for _, kw := range configs.V.Filter.Focuses {
		if strings.Contains(m.sraw, kw) {
			if flag == SPAM {
				flag = CONFUSION
			} else {
				flag = FOCUSED
			}
		}
	}
	m.tag = flag
	return m
}

// folder path should less then 240 bytes,
// file path should less than 260 bytes
func (m *Mail) deliver() error {
	if m.tag == DROP {
		return nil // Drop the mail, just ignore it
	}
	tag := func() string {
		switch m.tag {
		case FOCUSED:
			return "[FOCUSED]"
		case SPAM:
			return "[SPAM]"
		case CONFUSION:
			return "[CONFUSION]"
		default:
			return "[COMMON]"
		}
	}

	dstDir := filepath.Join(configs.V.MailSet, tag())
	if len(dstDir) >= 240 {
		return fmt.Errorf("Too long path: %s", dstDir)
	}
	if err := os.MkdirAll(dstDir, 0750); err != nil {
		return err
	}
	dstPath := filepath.Join(dstDir, filepath.Base(m.path))
	if len(dstPath) >= 260 {
		return fmt.Errorf("Too long file name: %s", dstPath)
	}
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, bytes.NewReader(m.raw)); err != nil {
		return err
	}
	return nil
}
