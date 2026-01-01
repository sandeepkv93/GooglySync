package fswatch

import "go.uber.org/zap"

// Watcher observes local filesystem changes.
type Watcher struct {
	Logger *zap.Logger
}

// NewWatcher constructs a filesystem watcher.
func NewWatcher(logger *zap.Logger) (*Watcher, error) {
	logger.Info("fswatch initialized")
	return &Watcher{Logger: logger}, nil
}
