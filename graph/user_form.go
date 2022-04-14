package graph

import (
	"auction-back/models"

	"github.com/hashicorp/go-multierror"
)

func isUserFormOwner(viewer models.User, form models.UserForm) error {
	if form.UserID != viewer.ID {
		return ErrViewerNotOwner
	}

	return nil
}

func (r *userFormResolver) isUserFormOwnerOrManager(viewer models.User, obj models.UserForm) error {
	var errors error

	if err := isUserFormOwner(viewer, obj); err != nil {
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
