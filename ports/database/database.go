package database

import (
	"gorm.io/gorm"
)

type Database struct {
	db gorm.DB
}

func New(db *gorm.DB) Database {
	return Database{db: *db}
}
