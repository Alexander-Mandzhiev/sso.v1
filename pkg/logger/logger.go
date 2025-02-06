package sl

import (
	"log/slog"
	"os"
	"sso/internal/config"
)

var Log *slog.Logger

const (
	envDevelopment = "development"
	envProduction  = "production"
)

func SetupLogger() {
	var level slog.Level
	switch config.Cfg.Env {
	case envDevelopment:
		level = slog.LevelDebug
	case envProduction:
		level = slog.LevelInfo
	default:
		level = slog.LevelInfo
	}
	Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
