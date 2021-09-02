package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hi20160616/spamkiller/configs"
)

func main() {
	defer configs.LogWriter.Close()
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	ms, err := InitializeMails(context.Background(), "spamkiller")
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
