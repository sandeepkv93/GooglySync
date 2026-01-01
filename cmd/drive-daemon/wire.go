//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/sandeepkv93/GooglySync/internal/auth"
	"github.com/sandeepkv93/GooglySync/internal/config"
	"github.com/sandeepkv93/GooglySync/internal/daemon"
	"github.com/sandeepkv93/GooglySync/internal/fswatch"
	"github.com/sandeepkv93/GooglySync/internal/logging"
	syncer "github.com/sandeepkv93/GooglySync/internal/sync"
	"github.com/sandeepkv93/GooglySync/internal/storage"
)

func InitializeDaemon() (*daemon.Daemon, error) {
	wire.Build(
		config.NewConfig,
		logging.NewLogger,
		storage.NewStorage,
		auth.NewService,
		fswatch.NewWatcher,
		syncer.NewEngine,
		daemon.NewDaemon,
	)
	return &daemon.Daemon{}, nil
}
