package models

import "gorm.io/gorm"

type Bank struct {
	gorm.Model
	ID   string `gorm:"default:generated();" json:"id"`
	Name string `json:"name"`
}
