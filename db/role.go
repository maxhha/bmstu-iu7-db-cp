package db

import (
	"database/sql"
	"time"
)

type RoleType string

const (
	RoleTypeManager RoleType = "MANAGER"
	RoleTypeAdmin   RoleType = "ADMIN"
)

type Role struct {
	Type      RoleType
	UserID    string
	User      User
	IssuerID  string
	Issuer    User
	CreatedAt time.Time `gorm:"default:now()"`
	DeletedAt sql.NullTime
}
