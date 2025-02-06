package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log/slog"
	"sso/internal/config"
	sl "sso/pkg/logger"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type Repository struct {
	db *sql.DB
}

func New(storagePath string) (*Repository, error) {
	const op = "repository.mssql.New"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("mssql", storagePath)
	if err != nil {
		sl.Log.Error(op, "Error opening connection to MSSQL database", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(config.Cfg.DBConfig.MaxOpenConnections)
	db.SetMaxIdleConns(config.Cfg.DBConfig.MaxIdleConnections)
	db.SetConnMaxLifetime(config.Cfg.DBConfig.ConnectionMaxLifetime)

	err = db.PingContext(ctx)
	if err != nil {
		sl.Log.Error(op, "Error testing MSSQL database connection", slog.Any("error", err))
		_ = db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sl.Log.Info(op, "Opened connection to MSSQL database")
	return &Repository{db: db}, nil
}

func (r *Repository) Stop() error {
	const op = "repository.mssql.Stop"

	if r.db == nil {
		sl.Log.Warn(op, "Database connection is already closed")
		return nil
	}

	err := r.db.Close()
	if err != nil {
		sl.Log.Error(op, "Error closing MSSQL database connection", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	r.db = nil
	sl.Log.Info(op, "MSSQL database connection closed")
	return nil
}
