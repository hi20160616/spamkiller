package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hi20160616/spamkiller/configs"
)

type Options struct {
	ProjectName configs.ProjectName
	MailsPath   MailsPath
}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := &Options{ProjectName: "spamkiller", MailsPath: MailsPath(os.Args[1])}
	app, err := NewApp(ctx, opts)
	if err != nil {
		app.ms.log.Fatal(err)
	}
	if err = app.work(ctx); err != nil {
		app.ms.log.Fatal(err)
	}
	if app.ms.cfg.Verbose {
		fmt.Println("Done. Press Enter to quit!")
		// bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Scanln()
	}
}

type app struct {
	ms   *Mails
	opts *Options
}

func NewApp(ctx context.Context, opts *Options) (*app, error) {
	ms, err := InitializeMails(ctx, opts.ProjectName, opts.MailsPath)
	if err != nil {
		return nil, err
	}
	return &app{ms, opts}, nil
}

func (a *app) work(ctx context.Context) error {
	a.ms.cfg.Folder = string(a.opts.MailsPath)

	if a.ms.cfg.Result == "./" {
		a.ms.cfg.Result = os.Args[1]
	}
	a.ms = a.ms.walkSrc(ctx)
	return a.ms.err
}
