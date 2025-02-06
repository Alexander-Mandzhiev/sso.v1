package models

import "github.com/google/uuid"

type UserInfo struct {
	Id    string    `json:"uid"`
	Name  string    `json:"name"`
	AppID int       `json:"app_id"`
	JTI   uuid.UUID `json:"jti"`
}
