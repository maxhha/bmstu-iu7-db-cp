package models

import (
	"gorm.io/gorm"
)

type ProductImage struct {
	gorm.Model
	ID        string `gorm:"type:varchar(16);"`
	Filename  string
	Path      string
	ProductID string `gorm:"type:varchar(16);"`
	Product   Product
}
