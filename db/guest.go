package db

import (
	"time"
)

type Guest struct {
	ID        string `gorm:"default:generated()"`
	ExpiresAt time.Time
}
