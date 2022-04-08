package database

import (
	"fmt"

	"gorm.io/gorm"
)

func paginationQueryByCreatedAtDesc(query *gorm.DB, first *int, after *string) (*gorm.DB, error) {
	pagination := query.Order("created_at desc")

	if first != nil {
		if *first < 1 {
			return nil, fmt.Errorf("first must be positive")
		}
		pagination = pagination.Limit(*first + 1)
	}

	if after != nil {
		afterCreatedAt := query.Where("id = ?", after).Select("created_at")
		pagination = pagination.Where("created_at < ANY(?)", afterCreatedAt)
	}

	return pagination, nil
}
