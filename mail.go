package main

import (
	"bytes"
	"io"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/hi20160616/spamkiller/configs"
)

const (
	COMMON int = iota
	FOCUSED
	SPAM
	CONFUSION
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
func (m *Mail) rename() error {
	return nil
}
