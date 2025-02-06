package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"sso/internal/config"
	"sso/internal/models"
	sl "sso/pkg/logger"
	"time"
)

func (r *Repository) CreateSession(ctx context.Context, userID string, appID int) (*models.Session, error) {
	const op = "repository.CreateSession"

	jti := uuid.New().String()
	expiresAt := time.Now().Add(config.Cfg.RefreshTokenTTL)

	query := `INSERT INTO sessions (jti, user_id, app_id, created_at, expires_at, is_active) 
		VALUES (?, ?, ?, ?, ?, 1)`

	result, err := r.db.ExecContext(ctx, query, jti, userID, appID, time.Now(), expiresAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		sl.Log.Warn(op, "Failed to create or update session", slog.String("user_id", userID))
		return nil, fmt.Errorf("%s: failed to create or update session", op)
	}

	session := &models.Session{
		JTI:       jti,
		UserID:    userID,
		AppID:     appID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsActive:  true,
	}

	sl.Log.Info(op, "Session created successfully", slog.String("jti", jti))
	return session, nil
}
