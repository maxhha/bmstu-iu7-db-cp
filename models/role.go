package models

import (
	"database/sql"
	"time"
)

// TODO: add role addition and removeing
type Role struct {
	Type      RoleType
	UserID    string
	User      User
	IssuerID  string
	Issuer    User
	CreatedAt time.Time `gorm:"default:now()"`
	DeletedAt sql.NullTime
}
