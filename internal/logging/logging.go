package logging

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/sandeepkv93/googlysync/internal/config"
)

// NewLogger builds a structured logger based on config.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	var ws zapcore.WriteSyncer
	if cfg.LogFilePath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.LogFilePath), 0o700); err != nil {
			return nil, err
		}
		lj := &lumberjack.Logger{
			Filename:   cfg.LogFilePath,
			MaxSize:    cfg.LogFileMaxMB,
			MaxBackups: cfg.LogFileMaxBackups,
			MaxAge:     cfg.LogFileMaxAgeDays,
			Compress:   true,
		}
		ws = zapcore.AddSync(lj)
	} else {
		ws = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, ws, level)
	return zap.New(core), nil
}
