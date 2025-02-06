package repository

import (
	"context"
	"fmt"
	"log/slog"
	sl "sso/pkg/logger"
)

func (r *Repository) DeactivateSession(ctx context.Context, jti string) error {
	const op = "repository.DeactivateSession"

	query := `UPDATE sessions SET is_active = 0 WHERE jti = ? AND is_active = 1`

	result, err := r.db.ExecContext(ctx, query, jti)
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

	sl.Log.Info(op, "Session deactivated successfully", slog.String("jti", jti))
	return nil
}
