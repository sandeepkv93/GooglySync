package auth

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/zalando/go-keyring"

	"github.com/sandeepkv93/googlysync/internal/config"
	"github.com/sandeepkv93/googlysync/internal/storage"
)

// State captures the current auth status.
type State struct {
	SignedIn bool
	Account  storage.Account
}

// Service handles auth and token lifecycle.
type Service struct {
	logger *zap.Logger
	cfg    *config.Config
	store  *storage.Storage
	krSvc  string

	mu    sync.Mutex
	state State
}

// NewService constructs the auth service.
func NewService(logger *zap.Logger, cfg *config.Config, store *storage.Storage) (*Service, error) {
	if logger == nil {
		return nil, errors.New("auth: logger is required")
	}
	if cfg == nil {
		return nil, errors.New("auth: config is required")
	}
	if store == nil {
		return nil, errors.New("auth: storage is required")
	}

	krSvc := cfg.AppName
	if krSvc == "" {
		krSvc = "googlysync"
	}
	svc := &Service{logger: logger, cfg: cfg, store: store, krSvc: krSvc}
	svc.bootstrapState(context.Background())
	logger.Info("auth service initialized")
	return svc, nil
}

// State returns the latest auth state.
func (s *Service) State() State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

// SignIn runs the OAuth flow and persists account metadata + refresh token.
func (s *Service) SignIn(ctx context.Context, scopes []string) error {
	if s.cfg.OAuthClientID == "" {
		return errors.New("oauth client id not configured")
	}
	if s.cfg.OAuthClientSecret == "" {
		return errors.New("oauth client secret not configured")
	}
	if len(scopes) == 0 {
		scopes = defaultScopes()
	}

	token, claims, err := runOAuthFlow(ctx, s.cfg, scopes, s.logger)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("oauth token missing")
	}

	accountID := claims.Sub
	if accountID == "" {
		return errors.New("oauth sub claim missing")
	}
	account := storage.Account{
		ID:          accountID,
		Email:       claims.Email,
		DisplayName: claims.Name,
		IsPrimary:   s.isFirstAccount(ctx),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.store.UpsertAccount(ctx, &account); err != nil {
		return err
	}

	refreshToken := token.RefreshToken
	if refreshToken == "" {
		return errors.New("refresh token missing; re-auth with consent")
	}
	ref := storage.TokenRef{
		AccountID: accountID,
		KeyID:     accountID,
		TokenType: "refresh",
		Scope:     scopeString(scopes),
		Expiry:    token.Expiry,
		UpdatedAt: time.Now(),
	}
	if err := s.store.UpsertTokenRef(ctx, &ref); err != nil {
		return err
	}
	if err := keyring.Set(s.krSvc, accountID, refreshToken); err != nil {
		_ = s.store.DeleteTokenRef(ctx, accountID)
		return err
	}

	s.mu.Lock()
	s.state = State{SignedIn: true, Account: account}
	s.mu.Unlock()

	return nil
}

// SignOut removes stored token reference and resets auth state.
func (s *Service) SignOut(ctx context.Context, accountID string) error {
	if accountID == "" {
		return errors.New("account id is required")
	}
	_ = keyring.Delete(s.krSvc, accountID)
	if err := s.store.DeleteAccount(ctx, accountID); err != nil {
		return err
	}
	s.mu.Lock()
	s.state = State{}
	s.mu.Unlock()
	return nil
}

func (s *Service) isFirstAccount(ctx context.Context) bool {
	accounts, err := s.store.ListAccounts(ctx)
	if err != nil {
		return false
	}
	return len(accounts) == 0
}

func (s *Service) bootstrapState(ctx context.Context) {
	accounts, err := s.store.ListAccounts(ctx)
	if err != nil || len(accounts) == 0 {
		return
	}

	var primary *storage.Account
	var fallback *storage.Account
	for i := range accounts {
		ref, err := s.store.GetTokenRef(ctx, accounts[i].ID)
		if err != nil || ref == nil {
			continue
		}
		if accounts[i].IsPrimary {
			primary = &accounts[i]
			break
		}
		if fallback == nil {
			fallback = &accounts[i]
		}
	}
	if primary == nil {
		primary = fallback
	}
	if primary == nil {
		return
	}

	s.mu.Lock()
	s.state = State{SignedIn: true, Account: *primary}
	s.mu.Unlock()
}

func scopeString(scopes []string) string {
	if len(scopes) == 0 {
		return ""
	}
	seen := make(map[string]struct{}, len(scopes))
	unique := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		unique = append(unique, scope)
	}
	if len(unique) == 0 {
		return ""
	}
	return strings.Join(unique, " ")
}
