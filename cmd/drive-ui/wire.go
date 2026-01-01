//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package main

import (
	"github.com/google/wire"

	"github.com/sandeepkv93/GooglySync/internal/config"
	"github.com/sandeepkv93/GooglySync/internal/logging"
)

func InitializeUI() (*UIApp, error) {
	wire.Build(
		config.NewConfig,
		logging.NewLogger,
		NewUIApp,
	)
	return &UIApp{}, nil
}
