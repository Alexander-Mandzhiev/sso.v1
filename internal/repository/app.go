package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sso/internal/models"
)

func (r *Repository) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.mssql.App"
	var app models.App

	if err := r.db.QueryRowContext(ctx, "SELECT id, name, secret FROM apps WHERE id = ?", id).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: failed to scan row: %w", op, err)
	}

	return app, nil
}
