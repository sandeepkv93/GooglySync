package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sandeepkv93/GooglySync/internal/config"
)

// NewLogger builds a structured logger based on config.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return nil, err
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = level
	zapCfg.Encoding = "json"
	zapCfg.EncoderConfig.TimeKey = "ts"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapCfg.Build()
}
