package role

import (
	"auction-back/auth"
	"auction-back/models"
	"auction-back/ports"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

type Interface interface {
	HasRole(roleType models.RoleType, viewer models.User) error
}

type RolePort struct {
	db ports.DB
}

func New(db ports.DB) RolePort {
	return RolePort{db}
}

func (r *RolePort) Handler() func(ctx context.Context, obj interface{}, next graphql.Resolver, role models.RoleType) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, role models.RoleType) (res interface{}, err error) {
		viewer, err := auth.ForViewer(ctx)

		if err != nil {
			return nil, err
		}

		if err := r.HasRole(role, viewer); err != nil {
			return nil, err
		}

		return next(ctx)
	}
}

func (r *RolePort) HasRole(roleType models.RoleType, viewer models.User) error {
	roles, err := r.db.Role().Find(ports.RoleFindConfig{
		Limit:   1,
		UserIDs: []string{viewer.ID},
		Types:   []models.RoleType{roleType},
	})
	if err != nil {
		return fmt.Errorf("db find: %w", err)
	}

	if len(roles) == 0 {
		return fmt.Errorf("no role %s", roleType)
	}

	return nil
}
