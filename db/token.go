package db

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
)

type TokenAction string

const (
	TokenActionApproveUserEmail TokenAction = "APPROVE_USER_EMAIL"
	TokenActionApproveUserPhone TokenAction = "APPROVE_USER_PHONE"
)

type Token struct {
	ID          uint
	UserID      string
	User        User
	CreatedAt   time.Time `gorm:"default:now();"`
	ActivatedAt sql.NullTime
	ExpiresAt   time.Time
	Action      TokenAction
	Data        datatypes.JSONMap
}
