package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string `json:"id"`
	CreatedAt    time.Time
	DeletedAt    sql.NullTime
	BlockedUntil sql.NullTime
}
