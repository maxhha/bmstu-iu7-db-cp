package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestConvertError(t *testing.T) {
	for err, target := range errorMap {
		assert.ErrorIs(t, convertError(err), target)
	}

	errors := []error{gorm.ErrInvalidDB, gorm.ErrInvalidValue, sql.ErrTxDone, sql.ErrConnDone}
	for _, err := range errors {
		assert.ErrorIs(t, convertError(err), err)
	}
}
