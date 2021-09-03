//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/hi20160616/spamkiller/configs"
)

func InitializeMails(ctx context.Context, projectName string, src MailsPath) (*Mails, error) {
	wire.Build(NewLog, NewMails, configs.NewConfig)
	return &Mails{}, nil
}