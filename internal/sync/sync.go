package sync

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/sandeepkv93/googlysync/internal/status"
	"github.com/sandeepkv93/googlysync/internal/storage"
)

// Engine coordinates sync operations.
type Engine struct {
	Logger *zap.Logger
	Store  *storage.Storage
	Status *status.Store
}

// NewEngine constructs a sync engine.
func NewEngine(logger *zap.Logger, store *storage.Storage, statusStore *status.Store) (*Engine, error) {
	logger.Info("sync engine initialized")
	return &Engine{Logger: logger, Store: store, Status: statusStore}, nil
}

// Run runs a stub sync loop that updates status periodically.
func (e *Engine) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.Status.Update(status.Snapshot{State: status.StateIdle, Message: "idle"})
			return
		case <-ticker.C:
			e.Status.Update(status.Snapshot{State: status.StateSyncing, Message: "sync tick"})
			e.Logger.Info("sync tick")
			e.Status.Update(status.Snapshot{State: status.StateIdle, Message: "idle"})
		}
	}
}
