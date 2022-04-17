package graph

import (
	"auction-back/models"
)

func isAccountOwner(viewer models.User, account models.Account) error {
	if account.UserID != viewer.ID {
		return ErrUserNotOwner
	}

	return nil
}
