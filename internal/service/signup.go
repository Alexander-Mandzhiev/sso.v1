package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sso/internal/config"
	"sso/internal/models"
	"sso/pkg/jwt"
	sl "sso/pkg/logger"
)

func (s *Service) SignUp(ctx context.Context, username string, password string, appID int) (string, string, error) {
	const op = "service.SignUp"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		sl.Log.Error(op, "failed hashing password", sl.Err(err))
		return "", "", err
	}

	// repository not implemented
	user, err := s.saveUser.SaveUser(ctx, username, passwordHash)
	if err != nil {
		return "", "", fmt.Errorf("failed getting data values in mssql database: %w", err)
	}

	app, err := s.appProvider.App(ctx, appID)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	userInfo := models.UserInfo{Id: user.ID, Name: user.Name, AppID: appID}

	session, err := s.sessionProvider.CreateSession(ctx, user.ID, appID)
	if err != nil {
		sl.Log.Error(op, "Failed to create session", sl.Err(err))
		return "", "", status.Error(codes.Internal, "failed to create session")
	}

	accessToken, err := jwt.GenerateToken(&userInfo, &app, config.Cfg.AccessTokenTTL, session.JTI)
	if err != nil {
		sl.Log.Warn("failed to generate token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, sl.Err(err))
	}

	refreshToken, err := jwt.GenerateToken(&userInfo, &app, config.Cfg.RefreshTokenTTL, session.JTI)
	if err != nil {
		sl.Log.Warn("failed to generate token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, sl.Err(err))
	}

	return accessToken, refreshToken, nil
}
