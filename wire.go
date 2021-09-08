//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/hi20160616/spamkiller/configs"
)

func InitializeMails(ctx context.Context, projectName configs.ProjectName, src MailsPath) (*Mails, error) {
	wire.Build(NewLog, NewMails, configs.NewConfig)
	return &Mails{}, nil
}

// func InitApp(ctx context.Context, opts *Options) (*app, error) {
//         wire.Build(appSet)
//         return nil, nil
// }
