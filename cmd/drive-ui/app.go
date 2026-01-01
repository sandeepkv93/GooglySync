package main

import (
	"go.uber.org/zap"

	"github.com/sandeepkv93/GooglySync/internal/config"
)

// UIApp is a placeholder for future GTK app wiring.
type UIApp struct {
	Config *config.Config
	Logger *zap.Logger
}

func NewUIApp(cfg *config.Config, logger *zap.Logger) (*UIApp, error) {
	return &UIApp{Config: cfg, Logger: logger}, nil
}
