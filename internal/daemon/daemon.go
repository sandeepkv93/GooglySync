package daemon

import (
	"go.uber.org/zap"

	"github.com/sandeepkv93/googlysync/internal/auth"
	"github.com/sandeepkv93/googlysync/internal/config"
	"github.com/sandeepkv93/googlysync/internal/fswatch"
	syncer "github.com/sandeepkv93/googlysync/internal/sync"
	"github.com/sandeepkv93/googlysync/internal/storage"
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

// Run starts the daemon loop.
func (d *Daemon) Run() error {
	d.Logger.Info("daemon running")
	return nil
}
