package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string
	Available float64 `sql:"type:decimal(10,2);"`
}

type Product struct {
	gorm.Model
	ID          string
	Name        string
	Description *string
	IsOnMarket  bool
	OwnerID     string
	Owner       User
}
