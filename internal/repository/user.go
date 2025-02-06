package repository

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sso/internal/models"
)

func (r *Repository) User(ctx context.Context, username string) (models.User, error) {
	const op = "repository.User"
	var user models.User
	err := r.db.QueryRowContext(ctx, `SELECT ID, Name FROM v8users WHERE Name = ?`, username).Scan(&user.ID, &user.Name)
	user.ID = hex.EncodeToString([]byte(user.ID))

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("%s: %s", op, ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
