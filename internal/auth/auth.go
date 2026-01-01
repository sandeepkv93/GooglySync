package auth

import (
	"go.uber.org/zap"

	"github.com/sandeepkv93/googlysync/internal/storage"
)

// Service handles auth and token lifecycle.
type Service struct {
	Logger *zap.Logger
	Store  *storage.Storage
}

// NewService constructs the auth service.
func NewService(logger *zap.Logger, store *storage.Storage) (*Service, error) {
	logger.Info("auth service initialized")
	return &Service{Logger: logger, Store: store}, nil
}
