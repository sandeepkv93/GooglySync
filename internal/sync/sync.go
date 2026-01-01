package sync

import (
	"go.uber.org/zap"

	"github.com/sandeepkv93/googlysync/internal/storage"
)

// Engine coordinates sync operations.
type Engine struct {
	Logger *zap.Logger
	Store  *storage.Storage
}

// NewEngine constructs a sync engine.
func NewEngine(logger *zap.Logger, store *storage.Storage) (*Engine, error) {
	logger.Info("sync engine initialized")
	return &Engine{Logger: logger, Store: store}, nil
}
