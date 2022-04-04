package db

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID          string `gorm:"default:generated();" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsOnMarket  bool   `gorm:"default:false;" json:"isOnMarket"`
	CreatorID   string
	Creator     User
}
