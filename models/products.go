package models

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          string       `json:"id"`
	State       ProductState `json:"state"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreatorID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}
