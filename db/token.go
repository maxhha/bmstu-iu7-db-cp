package db

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TokenAction string

const (
	TokenActionApproveUserEmail TokenAction = "APPROVE_USER_EMAIL"
	TokenActionApproveUserPhone TokenAction = "APPROVE_USER_PHONE"
)

type Token struct {
	gorm.Model
	ID          uint
	CreatorID   string
	Creator     TokenCreator
	ActivatedAt sql.NullTime
	ExpiresAt   time.Time
	Action      TokenAction
	Data        datatypes.JSONMap
}
