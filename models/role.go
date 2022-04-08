package models

import (
	"database/sql"
	"time"
)

// TODO: add role addition and removeing
type Role struct {
	Type      RoleType
	UserID    string
	IssuerID  string
	CreatedAt time.Time
	DeletedAt sql.NullTime
}
