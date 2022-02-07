package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string
	Available float64 `sql:"type:decimal(10,2);"`
}
