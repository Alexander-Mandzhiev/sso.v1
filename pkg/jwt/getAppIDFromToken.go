package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"sso/internal/models"
)

func GetAppIDFromToken(tokenString string) (int, error) {
	const op = "jwt.GetAppIDFromToken"

	parser := jwt.Parser{}

	var customClaims models.CustomClaims
	token, _, err := parser.ParseUnverified(tokenString, &customClaims)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to parse token without verification: %w", op, err)
	}

	if token == nil || token.Claims == nil {
		return 0, fmt.Errorf("%s: invalid token format", op)
	}

	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || claims.App == 0 {
		return 0, fmt.Errorf("%s: app_id is missing or invalid in token", op)
	}

	return claims.App, nil
}
