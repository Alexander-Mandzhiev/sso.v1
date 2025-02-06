package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"sso/internal/models"
	"sso/pkg/jwt"
	sl "sso/pkg/logger"
)

func (s *Service) app(ctx context.Context, refreshToken string) (models.App, error) {
	op := "service.UserInfo"
	appID, err := jwt.GetAppIDFromToken(refreshToken)
	if err != nil {
		sl.Log.Error(op, "Failed to get app_id from token", slog.Any("error", err))
		return models.App{}, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	if appID == 0 {
		sl.Log.Warn(op, "Invalid app_id in token")
		return models.App{}, status.Error(codes.Unauthenticated, "invalid app_id in token")
	}

	app, err := s.appProvider.App(ctx, appID)
	if err != nil {
		sl.Log.Error(op, "Failed to fetch app secret", slog.Int("app_id", appID), slog.Any("error", err))
		return models.App{}, status.Error(codes.Internal, "failed to fetch app info")
	}

	return app, nil
}
