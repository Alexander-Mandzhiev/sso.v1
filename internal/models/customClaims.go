package models

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	JTI    string `json:"jti"`
	App    int    `json:"app"`
	Issuer string `json:"iss"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
