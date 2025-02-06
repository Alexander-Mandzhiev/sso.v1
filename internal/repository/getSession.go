package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/models"
	sl "sso/pkg/logger"
)

func (r *Repository) GetSession(ctx context.Context, jti string) (*models.Session, error) {
	const op = "repository.GetSession"

	query := `SELECT jti, user_id, app_id, created_at, expires_at, is_active FROM sessions WHERE jti = ?`

	var session models.Session
	err := r.db.QueryRowContext(ctx, query, jti).Scan(&session.JTI, &session.UserID, &session.AppID, &session.CreatedAt, &session.ExpiresAt, &session.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sl.Log.Warn(op, "Session not found", slog.String("jti", jti))
			return nil, fmt.Errorf("%s: session not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sl.Log.Info(op, "Session fetched successfully", slog.String("jti", jti))
	return &session, nil
}
