package service

import (
	"context"
	"errors"
	"sso/internal/models"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserProvider interface {
	User(ctx context.Context, username string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type SessionProvider interface {
	CreateSession(ctx context.Context, userID string, appID int) (*models.Session, error)
	DeactivateSession(ctx context.Context, jti string) error
	GetSession(ctx context.Context, jti string) (*models.Session, error)
	UpdateSession(ctx context.Context, jti string, ttl time.Duration) error
}

type Service struct {
	userProvider    UserProvider
	appProvider     AppProvider
	sessionProvider SessionProvider
}

func New(userProvider UserProvider, appProvider AppProvider, sessionProvider SessionProvider) *Service {
	return &Service{userProvider: userProvider, appProvider: appProvider, sessionProvider: sessionProvider}
}
