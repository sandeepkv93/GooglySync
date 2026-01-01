package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"

	"github.com/sandeepkv93/googlysync/internal/config"
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

// Close shuts down the database connection.
func (s *Storage) Close() error {
	if s == nil || s.DB == nil {
		return nil
	}
	return s.DB.Close()
}
