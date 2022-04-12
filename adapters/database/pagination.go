package database

import (
	"fmt"

	"gorm.io/gorm"
)

func paginationQueryByCreatedAtDesc(query *gorm.DB, first *int, after *string) (*gorm.DB, error) {
	if after != nil {
		afterCreatedAt := query.
			Session(&gorm.Session{Initialized: true}).
			Model(query.Statement.Model).
			Where("id = ?", after).
			Select("created_at")

		query = query.Where("created_at < ANY( ? )", afterCreatedAt)
	}

	if first != nil {
		if *first < 1 {
			return nil, fmt.Errorf("first must be positive")
		}
		query = query.Limit(*first + 1)
	}

	query = query.Order("created_at desc")

	return query, nil
}
