package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hi20160616/spamkiller/configs"
	"github.com/pkg/errors"
)

const (
	ERR int = iota
	COMMON
	FOCUSED
	SPAM
	CONFUSION
	DROP
)

type Mails struct {
	cfg   *configs.Config
	src   string // mails folder
	mails []*Mail
	err   error
	log   *log.Logger
}

func NewMails(ctx context.Context, cfg *configs.Config, log *log.Logger, src string) *Mails {
	return &Mails{
		cfg: cfg,
		log: log,
		src: src,
	}
}

func (ms *Mails) walkSrc(ctx context.Context) *Mails {
	subDirToSkip := "skip"
	ms.err = filepath.WalkDir(
		ms.src, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				ms.log.Printf(
					"prevent panic by handling failure accessing a path %q: %v\n",
					path, err)
				return err
			}
			if d.IsDir() && d.Name() == subDirToSkip {
				ms.log.Printf(
					"skipping a dir without errors: %+v \n", d.Name())
				return filepath.SkipDir
			}
			if d.Type().IsRegular() {
				if err := treat(ctx, ms.cfg, ms.log, path); err != nil {
					ms.log.Println(err)
					// return err
				}
			}
			return nil
		})
	return ms
}

func trueType(cfg *configs.Config, src string) error {
	ext := filepath.Ext(src)
	trueType := func() bool {
		for _, e := range cfg.MailTypes {
			if e == ext {
				return true
			}
		}
		return false
	}()
	if !trueType {
		return errors.WithMessagef(ErrNotSupportExt,
			"bypass this file: %s", src)
	}
	return nil
}

func treat(ctx context.Context, cfg *configs.Config, log *log.Logger, src string) error {
	if err := trueType(cfg, src); err != nil {
		return err
	}
	// new mail
	m := NewMail(ctx, cfg, log, src)
	if m.err != nil {
		// if m.log != nil {
		//         m.log.Println(m.err)
		// }
		// log.Println(m.err)
		m.log.Println(m.err)
	}
	// analysis and deliver, log out the error
	return m.analysis().deliver()
}

type Mail struct {
	cfg           *configs.Config
	log           *log.Logger
	src           string // source path
	raw           []byte
	sraw          string
	msg           *mail.Message
	date          time.Time
	from          *mail.Address
	to            []*mail.Address
	subject, body string
	tag           int
	err           error
}

func NewMail(ctx context.Context, cfg *configs.Config, log *log.Logger, mailPath string) *Mail {
	rtErr := func(err error) *Mail {
		return &Mail{src: mailPath, cfg: cfg, log: log, err: err}
	}
	r, err := os.ReadFile(mailPath)
	if err != nil {
		return rtErr(errors.WithMessagef(err, "ReadFile error: %s", mailPath))
	}
	msg, err := mail.ReadMessage(bytes.NewReader(r))
	if err != nil {
		return rtErr(errors.WithMessagef(err, "ReadMessage error: %s", mailPath))
	}
	subject := msg.Header.Get("Subject")
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return rtErr(errors.WithMessagef(err, "Read body error: %s", mailPath))
	}
	date, err := msg.Header.Date()
	if err != nil {
		return rtErr(errors.WithMessagef(err, "Read Date error: %s", mailPath))
	}
	from, err := msg.Header.AddressList("From")
	if err != nil {
		return rtErr(errors.WithMessagef(err, "Read Address From error: %s", mailPath))
	}
	to, err := msg.Header.AddressList("To")
	if err != nil {
		return rtErr(errors.WithMessagef(err, "Read Address To error: %s", mailPath))
	}
	sraw := msg.Header.Get("From") + msg.Header.Get("To") +
		subject + " " + string(body)
	return &Mail{
		cfg:     cfg,
		log:     log,
		src:     mailPath,
		raw:     r,
		sraw:    sraw,
		msg:     msg,
		date:    date,
		from:    from[0],
		to:      to,
		subject: subject,
		body:    string(body),
	}
}

func (m *Mail) analysis() *Mail {
	flag := COMMON
	if m.err != nil {
		m.tag = ERR
		return m
	}
	if m.date.Before(m.cfg.Drop) {
		m.tag = DROP
		return m
	}
	for _, kw := range m.cfg.Filter.Spams {
		if strings.Contains(m.sraw, kw) {
			flag = SPAM
		}
	}
	for _, kw := range m.cfg.Filter.Focuses {
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
		case ERR:
			return "[ERROR]"
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

	dstDir := filepath.Join(m.cfg.Result, tag())
	if len(dstDir) >= 240 {
		return fmt.Errorf("Too long path: %s", dstDir)
	}
	if err := os.MkdirAll(dstDir, 0750); err != nil {
		return err
	}
	dstPath := filepath.Join(dstDir, filepath.Base(m.src))
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
