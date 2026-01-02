package storage

import (
	"context"
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
	"go.uber.org/zap"

	"github.com/pressly/goose/v3"

	"github.com/sandeepkv93/googlysync/internal/config"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := migrate(context.Background(), db, logger); err != nil {
		_ = db.Close()
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

func migrate(ctx context.Context, db *sql.DB, logger *zap.Logger) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	logger.Info("applying migrations")
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return err
	}
	logger.Info("migrations complete")
	return nil
}
