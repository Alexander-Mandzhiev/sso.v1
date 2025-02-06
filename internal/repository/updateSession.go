package repository

import (
	"context"
	"fmt"
	"log/slog"
	sl "sso/pkg/logger"
	"time"
)

func (r *Repository) UpdateSession(ctx context.Context, jti string, ttl time.Duration) error {
	const op = "repository.UpdateSession"

	fmt.Println(int(ttl.Seconds()))

	query := `UPDATE sessions SET expires_at = DATEADD(SECOND, ?, GETDATE()) WHERE jti = ? AND is_active = 1`
	result, err := r.db.ExecContext(ctx, query, int(ttl.Seconds()), jti)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		sl.Log.Warn(op, "Session not found or already inactive", slog.String("jti", jti))
		return fmt.Errorf("%s: session not found or already inactive", op)
	}

	sl.Log.Info(op, "Session updated successfully", slog.String("jti", jti))
	return nil
}
