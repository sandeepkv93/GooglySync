package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
	"go.uber.org/zap"

	"github.com/sandeepkv93/GooglySync/internal/config"
)

// Storage wraps access to the local metadata store.
type Storage struct {
	DB *sql.DB
}

// NewStorage opens the SQLite database for metadata.
func NewStorage(cfg *config.Config, logger *zap.Logger) (*Storage, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0o700); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	logger.Info("storage initialized", zap.String("path", cfg.DatabasePath))
	return &Storage{DB: db}, nil
}
