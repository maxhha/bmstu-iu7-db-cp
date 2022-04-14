package graph

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func isProductOwner(DB ports.DB, viewer models.User, product models.Product) error {
	owner, err := DB.Product().GetOwner(product)
	if err != nil {
		return fmt.Errorf("get owner: %w", err)
	}

	if owner.ID != viewer.ID {
		return ErrViewerNotOwner
	}

	return nil
}

func (r *productResolver) isProductOwnerOrManager(viewer models.User, obj models.Product) error {
	var errors error

	if err := isProductOwner(r.DB, viewer, obj); err != nil {
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
