package database

import (
	"auction-back/ports"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var errorMap = map[error]error{
	gorm.ErrRecordNotFound: ports.ErrRecordNotFound,
}

func convertError(err error) error {
	for check, converted := range errorMap {
		if errors.Is(err, check) {
			return fmt.Errorf("%w (converted from: %v)", converted, err)
		}
	}

	return err
}
