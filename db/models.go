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

type Offer struct {
	gorm.Model
	ID         string
	Amount     float64 `sql:"type:decimal(10,2);"`
	ConsumerID string
	Consumer   User
	ProductID  string
	Product    Product
}
