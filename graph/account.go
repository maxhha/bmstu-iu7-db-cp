package graph

import (
	"auction-back/models"
)

func IsAccountOwner(viewer models.User, account models.Account) error {
	if account.UserID != viewer.ID {
		return ErrUserNotOwner
	}

	return nil
}
