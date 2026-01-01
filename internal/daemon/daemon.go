package daemon

import (
	"context"

	"go.uber.org/zap"

	"github.com/sandeepkv93/googlysync/internal/auth"
	"github.com/sandeepkv93/googlysync/internal/config"
	"github.com/sandeepkv93/googlysync/internal/fswatch"
	"github.com/sandeepkv93/googlysync/internal/storage"
	syncer "github.com/sandeepkv93/googlysync/internal/sync"
)

// Daemon wires together core services.
type Daemon struct {
	Logger  *zap.Logger
	Config  *config.Config
	Storage *storage.Storage
	Auth    *auth.Service
	Sync    *syncer.Engine
	Watcher *fswatch.Watcher
}

// NewDaemon constructs a daemon.
func NewDaemon(
	logger *zap.Logger,
	cfg *config.Config,
	store *storage.Storage,
	authSvc *auth.Service,
	syncEngine *syncer.Engine,
	watcher *fswatch.Watcher,
) (*Daemon, error) {
	logger.Info("daemon initialized")
	return &Daemon{
		Logger:  logger,
		Config:  cfg,
		Storage: store,
		Auth:    authSvc,
		Sync:    syncEngine,
		Watcher: watcher,
	}, nil
}

// Run starts the daemon loop and blocks until shutdown.
func (d *Daemon) Run(ctx context.Context) error {
	d.Logger.Info("daemon running")

	<-ctx.Done()
	d.Logger.Info("daemon shutting down")

	return d.Close()
}

// Close releases resources owned by the daemon.
func (d *Daemon) Close() error {
	if d.Storage != nil {
		return d.Storage.Close()
	}
	return nil
}
