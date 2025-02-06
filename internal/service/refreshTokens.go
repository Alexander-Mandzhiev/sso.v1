package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"sso/internal/config"
	"sso/pkg/jwt"
	sl "sso/pkg/logger"
	"time"
)

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "service.RefreshTokens"

	app, err := s.app(ctx, refreshToken)
	if err != nil {
		sl.Log.Error(op, "Failed to fetch app secret", slog.Any("error", err))
		return "", "", status.Error(codes.Internal, "failed to fetch app info")
	}

	userInfo, err := jwt.ValidateRefreshToken(ctx, refreshToken, app.Secret)
	if err != nil {
		sl.Log.Error(op, "Invalid refresh token", slog.String("refresh_token", refreshToken), slog.Any("error", err))
		return "", "", status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	session, err := s.sessionProvider.GetSession(ctx, userInfo.JTI.String())
	if err != nil {
		sl.Log.Error(op, "Failed to get session", sl.Err(err))
		return "", "", status.Error(codes.Internal, "failed to get session")
	}

	if !session.IsActive || session.ExpiresAt.Before(time.Now()) {
		sl.Log.Warn(op, "Session is inactive or expired", slog.String("jti", session.JTI))
		return "", "", status.Error(codes.Unauthenticated, "session is inactive or expired")
	}

	if err = s.sessionProvider.UpdateSession(ctx, session.JTI, config.Cfg.SSO.RefreshTokenTTL); err != nil {
		sl.Log.Error(op, "Failed to update session", sl.Err(err))
		return "", "", status.Error(codes.Internal, "failed to update session")
	}

	newAccessToken, err := jwt.GenerateToken(&userInfo, &app, config.Cfg.AccessTokenTTL, session.JTI)
	if err != nil {
		sl.Log.Error(op, "Failed to generate new access token", slog.Any("error", err))
		return "", "", status.Error(codes.Internal, "failed to generate access token")
	}

	newRefreshToken, err := jwt.GenerateToken(&userInfo, &app, config.Cfg.RefreshTokenTTL, session.JTI)
	if err != nil {
		sl.Log.Error(op, "Failed to generate new refresh token", slog.Any("error", err))
		return "", "", status.Error(codes.Internal, "failed to generate refresh token")
	}

	sl.Log.Info(op, "Tokens refreshed successfully", slog.String("user_id", userInfo.Id))
	return newAccessToken, newRefreshToken, nil
}
