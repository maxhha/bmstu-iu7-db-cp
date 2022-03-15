package db

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           string
	Email        string
	Phone        string
	Password     string
	Name         string
	BlockedUntil sql.NullTime
}
