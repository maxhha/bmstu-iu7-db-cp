package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string    `gorm:"default:generated();" json:"id"`
	CreatedAt    time.Time `gorm:"default:now();"`
	DeletedAt    sql.NullTime
	BlockedUntil sql.NullTime
}
