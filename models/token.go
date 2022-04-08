package models

import (
	"database/sql"
	"time"
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
	CreatedAt   time.Time `gorm:"default:now();"`
	ActivatedAt sql.NullTime
	ExpiresAt   time.Time
	Action      TokenAction
	Data        map[string]interface{}
}
