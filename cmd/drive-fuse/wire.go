//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package main

import (
	"github.com/google/wire"

	"github.com/sandeepkv93/googlysync/internal/config"
	"github.com/sandeepkv93/googlysync/internal/logging"
)

func InitializeFUSE() (*FUSEApp, error) {
	wire.Build(
		config.NewConfig,
		logging.NewLogger,
		NewFUSEApp,
	)
	return &FUSEApp{}, nil
}
