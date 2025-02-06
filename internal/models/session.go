package models

import "time"

type Session struct {
	JTI       string    `json:"jti"`
	UserID    string    `json:"user_id"`
	AppID     int       `json:"app_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}
