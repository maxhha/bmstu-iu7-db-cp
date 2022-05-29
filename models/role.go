package models

import (
	"database/sql"
	"time"
)

type RoleType string

var (
	RoleTypeManager RoleType = "MANAGER"
	RoleTypeAdmin   RoleType = "ADMIN"
)

// TODO: add role addition and removeing
type Role struct {
	Type      RoleType
	UserID    string
	IssuerID  string
	CreatedAt time.Time
	DeletedAt sql.NullTime
}
