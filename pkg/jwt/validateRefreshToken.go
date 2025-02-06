package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"sso/internal/models"
	sl "sso/pkg/logger"
)

func ValidateRefreshToken(ctx context.Context, token, secret string) (models.UserInfo, error) {
	const op = "jwt.ValidateRefreshToken"

	var customClaims models.CustomClaims
	_, err := jwt.ParseWithClaims(token, &customClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method: %v", op, token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		sl.Log.ErrorContext(ctx, op, "Failed to parse token", slog.Any("error", err))
		return models.UserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	userInfo := models.UserInfo{
		Id:    customClaims.Subject,
		Name:  customClaims.Name,
		AppID: customClaims.App,
		JTI:   uuid.MustParse(customClaims.JTI),
	}

	return userInfo, nil
}
