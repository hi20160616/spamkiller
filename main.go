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

	// opts := &Options{ProjectName: "spamkiller", MailsPath: MailsPath(os.Args[1])}
	ms, err := InitializeMails(context.Background(), "spamkiller", MailsPath(os.Args[1]))
	if err != nil {
		fmt.Println(err)
	}
	ms.cfg.Folder = os.Args[1]
	if ms.cfg.Result == "./" {
		ms.cfg.Result = os.Args[1]
	}
	ms = ms.walkSrc(context.Background())
	if ms.err != nil {
		fmt.Println(ms.err)
		fmt.Println("Sth error. Press Enter to quit!")
		// bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Scanln()

	}
	if ms.cfg.Verbose {
		fmt.Println("Done. Press Enter to quit!")
		// bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Scanln()
	}
}

// var appSet = wire.NewSet(
//         wire.Struct(new(app), "*"),
//         NewLog, NewMails, configs.NewConfig,
// )
//
// type app struct {
//         ms   *Mails
//         opts *Options
// }
//
// func (a *app) work(ctx context.Context) error {
//         a.ms.cfg.Folder = string(a.opts.MailsPath)
//
//         if a.ms.cfg.Result == "./" {
//                 a.ms.cfg.Result = os.Args[1]
//         }
//         a.ms = a.ms.walkSrc(ctx)
//         if a.ms.err != nil {
//                 fmt.Println(a.ms.err)
//                 fmt.Println("Sth error. Press Enter to quit!")
//                 // bufio.NewReader(os.Stdin).ReadBytes('\n')
//                 fmt.Scanln()
//                 return a.ms.err
//
//         }
//         if a.ms.cfg.Verbose {
//                 fmt.Println("Done. Press Enter to quit!")
//                 // bufio.NewReader(os.Stdin).ReadBytes('\n')
//                 fmt.Scanln()
//         }
//         return nil
// }
