package main

import (
	"go.uber.org/zap"

	"github.com/sandeepkv93/googlysync/internal/config"
)

// FUSEApp is a placeholder for future streaming mode wiring.
type FUSEApp struct {
	Config *config.Config
	Logger *zap.Logger
}

func NewFUSEApp(cfg *config.Config, logger *zap.Logger) (*FUSEApp, error) {
	return &FUSEApp{Config: cfg, Logger: logger}, nil
}
