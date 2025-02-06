package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"sso/pkg/jwt"
	sl "sso/pkg/logger"
)

func (s *Service) RevokeTokens(ctx context.Context, refreshToken string) error {
	const op = "service.RevokeTokens"

	app, err := s.app(ctx, refreshToken)
	if err != nil {
		sl.Log.Error(op, "Failed to fetch app secret", slog.Any("error", err))
		return status.Error(codes.Internal, "failed to fetch app info")
	}

	userInfo, err := jwt.ValidateRefreshToken(ctx, refreshToken, app.Secret)
	if err != nil {
		sl.Log.Error(op, "Invalid refresh token", sl.Err(err))
		return fmt.Errorf("%s: invalid refresh token: %w", op, err)
	}

	if err = s.sessionProvider.DeactivateSession(ctx, userInfo.JTI.String()); err != nil {
		sl.Log.Error(op, "Failed to deactivate session", sl.Err(err))
		return fmt.Errorf("%s: failed to deactivate session: %w", op, err)
	}

	sl.Log.Info(op, "Tokens revoked successfully", slog.String("user_id", userInfo.Id))
	return nil
}
