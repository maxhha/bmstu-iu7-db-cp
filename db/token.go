package db

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
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

func (t *Token) EnsureFillUser(db *gorm.DB) error {
	if t.UserID == t.User.ID {
		return nil
	}

	return db.Take(&t.User, "id = ?", t.UserID).Error
}
