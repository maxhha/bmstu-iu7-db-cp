package models

import (
	"database/sql"
	"time"
)

type Bank struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
