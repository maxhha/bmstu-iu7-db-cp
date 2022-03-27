package role

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/model"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"gorm.io/gorm"
)

type Interface interface {
	HasRole(roleType db.RoleType, viewer *db.User) error
}

type RolePort struct {
	db *gorm.DB
}

func New(db *gorm.DB) RolePort {
	return RolePort{db}
}

func (r *RolePort) Handler() func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
		if err := r.HasRole(db.RoleType(role), auth.ForViewer(ctx)); err != nil {
			return nil, err
		}

		return next(ctx)
	}
}

func (r *RolePort) HasRole(roleType db.RoleType, viewer *db.User) error {
	if viewer == nil {
		return fmt.Errorf("unauthorized")
	}

	role := make([]db.Role, 1)

	err := r.db.Limit(1).Find(&role, "user_id = ? AND type = ?", viewer.ID, roleType).Error
	if err != nil {
		return fmt.Errorf("find: %w", err)
	}

	if len(role) == 0 {
		return fmt.Errorf("no role %s", roleType)
	}

	return nil
}
