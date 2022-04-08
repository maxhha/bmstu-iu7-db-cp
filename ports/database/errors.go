package database

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

func convertError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return sql.ErrNoRows
	}
	return err
}
