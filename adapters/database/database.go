package database

import (
	"auction-back/ports"

	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func New(db *gorm.DB) Database {
	return Database{db: db}
}

func (d *Database) Tx() ports.TXDB {
	return &Database{d.db.Begin()}
}

func (d *Database) DB() ports.DB {
	return d
}

func (d *Database) Rollback() {
	d.db.Rollback()
}

func (d *Database) Commit() error {
	return d.db.Commit().Error
}
