package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string
	Available uint
}
