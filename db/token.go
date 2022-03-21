package db

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
)

type TokenAction string

const (
	TokenActionSetUserEmail     TokenAction = "SET_USER_EMAIL"
	TokenActionSetUserPhone     TokenAction = "SET_USER_PHONE"
	TokenActionModerateUserForm TokenAction = "MODERATE_USER_FORM"
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
