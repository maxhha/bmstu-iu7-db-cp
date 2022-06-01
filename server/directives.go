package server

import (
	"auction-back/auth"
	"auction-back/models"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/hashicorp/go-multierror"
)

func hasRoleDirective(roleChecker map[models.RoleEnum]RoleChecker) func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []models.RoleEnum) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []models.RoleEnum) (res interface{}, err error) {
		viewer, err := auth.ForViewer(ctx)
		if err != nil {
			return nil, err
		}

		var errors error
		for _, role := range roles {
			checker, exists := roleChecker[role]
			if !exists {
				return nil, fmt.Errorf("%w: %s", ErrUnknownRole, role)
			}

			if err := checker(viewer, obj); err != nil {
				errors = multierror.Append(errors, err)
			} else {
				return next(ctx)
			}
		}

		return nil, fmt.Errorf("denied: %w", errors)
	}
}
