package repository

import (
	"context"
	"sso/internal/models"
	sl "sso/pkg/logger"
)

func (r *Repository) SaveUser(ctx context.Context, username string, passwordHash []byte) (models.User, error) {
	const op = "storage.mssql.App"
	var user models.User

	/*
		if err := r.db.QueryRowContext(ctx, "SELECT id, name, secret FROM apps WHERE id = ?", id).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.User{}, fmt.Errorf("%s: %w", op, ErrAppNotFound)
			}
			return models.User{}, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}*/

	sl.Log.Info(op, "Save user successfully")
	return user, nil
}
