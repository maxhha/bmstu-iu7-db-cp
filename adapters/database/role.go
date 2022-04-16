package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Type      models.RoleType
	UserID    string
	IssuerID  string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (r *Role) into() models.Role {
	return models.Role{
		Type:      r.Type,
		UserID:    r.UserID,
		IssuerID:  r.IssuerID,
		CreatedAt: r.CreatedAt,
		DeletedAt: sql.NullTime(r.DeletedAt),
	}
}

type roleDB struct{ *Database }

func (d *Database) Role() ports.RoleDB { return &roleDB{d} }

func (d *roleDB) Find(config ports.RoleFindConfig) ([]models.Role, error) {
	query := d.db

	if config.Limit > 0 {
		query = query.Limit(config.Limit)
	}

	if len(config.Types) > 0 {
		query = query.Where("type IN ?", config.Types)
	}

	if len(config.UserIDs) > 0 {
		query = query.Where("user_id IN ?", config.UserIDs)
	}

	var objs []Role
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", convertError(err))
	}

	arr := make([]models.Role, 0, len(objs))
	for _, obj := range objs {
		arr = append(arr, obj.into())
	}

	return arr, nil
}
