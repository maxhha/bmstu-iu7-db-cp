package graph

import (
	"auction-back/models"
)

func IsUserFormOwner(viewer models.User, form models.UserForm) error {
	if form.UserID != viewer.ID {
		return ErrViewerNotOwner
	}

	return nil
}
