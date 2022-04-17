package models

import (
	"database/sql"
	"time"
)

type Account struct {
	ID               string
	Number           string
	UserID           string
	NominalAccountID string
	CreatedAt        time.Time
	DeletedAt        sql.NullTime
}
