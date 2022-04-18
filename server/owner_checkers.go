package server

import (
	"auction-back/graph"
	"auction-back/models"
	"errors"
	"fmt"
)

type OwnerChecker func(viewer models.User, obj interface{}) error

var ErrFailConvert = errors.New("fail convert")

func userOwnerChecker(viewer models.User, obj interface{}) error {
	user, ok := obj.(*models.User)
	if !ok {
		return fmt.Errorf("%w to User", ErrFailConvert)
	}

	return graph.IsUserOwner(viewer, *user)
}

func userFormOwnerChecker(viewer models.User, obj interface{}) error {
	userForm, ok := obj.(*models.UserForm)
	if !ok {
		return fmt.Errorf("%w to UserForm", ErrFailConvert)
	}

	return graph.IsUserFormOwner(viewer, *userForm)
}

func userFormFilledOwnerChecker(viewer models.User, obj interface{}) error {
	userForm, ok := obj.(*models.UserFormFilled)
	if !ok {
		return fmt.Errorf("%w to UserFormFilled", ErrFailConvert)
	}

	return graph.IsUserFormOwner(viewer, *userForm.UserForm)
}

func newOwnerCheckers() map[string]OwnerChecker {
	return map[string]OwnerChecker{
		"User":           userOwnerChecker,
		"UserForm":       userFormOwnerChecker,
		"UserFormFilled": userFormFilledOwnerChecker,
	}
}
