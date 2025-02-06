package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"sso/internal/models"
	"time"
)

func GenerateToken(userInfo *models.UserInfo, app *models.App, ttl time.Duration, jti string) (string, error) {
	claims := models.CustomClaims{
		JTI:    jti,
		App:    userInfo.AppID,
		Issuer: app.Name,
		Name:   userInfo.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    app.Name,
			Subject:   userInfo.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
