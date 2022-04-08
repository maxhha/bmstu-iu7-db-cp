package graph

import (
	"auction-back/models"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func (r *productResolver) isOwner(viewer models.User, product models.Product) error {
	owner, err := r.DB.Product().GetOwner(product)
	if err != nil {
		return fmt.Errorf("get owner: %w", err)
	}

	if owner.ID != viewer.ID {
		return fmt.Errorf("viewer is not owner")
	}

	return nil
}

func (r *productResolver) isOwnerOrManager(viewer models.User, obj models.Product) error {
	var errors error

	if err := r.isOwner(viewer, obj); err != nil {
		errors = multierror.Append(errors, err)
	} else {
		return nil
	}

	if err := r.RolePort.HasRole(models.RoleTypeManager, viewer); err != nil {
		errors = multierror.Append(errors, err)
	} else {
		return nil
	}

	return errors
}
