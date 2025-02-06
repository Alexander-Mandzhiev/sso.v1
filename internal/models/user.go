package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id,omitempty" db:"ID"`
	Name         string    `json:"username,omitempty" db:"Name"`
	FullName     string    `json:"full_name,omitempty" db:"Descr"`
	OSName       string    `json:"os_name,omitempty" db:"OSName"`
	Changed      time.Time `json:"created_at,omitempty" db:"Changed"`
	RolesID      int       `json:"active,omitempty" db:"RolesID"`
	Show         bool      `json:"show,omitempty" db:"Show"`
	PasswordHash []byte    `json:"password_hash,omitempty" db:"Data"`
	EAuth        bool      `json:"e_auth,omitempty" db:"EAuth"`
	Admin        bool      `json:"admin,omitempty" db:"AdmRole"`
	UserSprH     int       `json:"us_spr_h,omitempty" db:"UsSprH"`
}
