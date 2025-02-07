package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"sso/pkg/jwt"
	sl "sso/pkg/logger"
	"time"
)

func (s *Service) ValidateToken(ctx context.Context, accessToken string) (bool, string, int, error) {
	const op = "service.ValidateToken"

	app, err := s.app(ctx, accessToken)
	if err != nil {
		sl.Log.Error(op, "Failed to fetch app secret", slog.Any("error", err))
		return false, "", 0, status.Error(codes.Internal, "failed to fetch app info")
	}

	userInfo, err := jwt.ValidateRefreshToken(ctx, accessToken, app.Secret)
	if err != nil {
		sl.Log.Error(op, "Invalid refresh token", slog.String("refresh_token", accessToken), slog.Any("error", err))
		return false, "", 0, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	user, err := s.userProvider.User(ctx, userInfo.Name)
	if err != nil {
		return false, "", 0, fmt.Errorf("failed getting data values in mssql database: %w", err)
	}

	session, err := s.sessionProvider.GetSession(ctx, userInfo.JTI.String())
	if err != nil {
		sl.Log.Error(op, "Failed to get session", sl.Err(err))
		return false, "", 0, status.Error(codes.Internal, "failed to get session")
	}

	if !session.IsActive || session.ExpiresAt.Before(time.Now()) {
		sl.Log.Warn(op, "Session is inactive or expired", slog.String("jti", session.JTI))
		return false, "", 0, status.Error(codes.Unauthenticated, "session is inactive or expired")
	}

	return true, user.Name, app.ID, nil
}
