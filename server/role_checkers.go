package server

import (
	"auction-back/models"
	"auction-back/ports"
	"errors"
	"fmt"
	"reflect"
)

type RoleChecker func(viewer models.User, obj interface{}) error

var ErrUnknownRole = errors.New("unknown role")
var ErrUnexpectObjectType = errors.New("unexpected object type")

func userRoleChecker(viewer models.User, obj interface{}) error {
	return nil
}

func roleTypeRoleChecker(r ports.Role, role models.RoleType) RoleChecker {
	return func(viewer models.User, obj interface{}) error {
		return r.HasRole(role, viewer)
	}
}

func ownerRoleChecker(ownerChecker map[string]OwnerChecker) RoleChecker {
	return func(viewer models.User, obj interface{}) error {
		t := reflect.TypeOf(obj)

		if t.Kind() != reflect.Ptr {
			return fmt.Errorf("%w: obj is not ptr", ErrUnexpectObjectType)
		}

		name := t.Elem().Name()
		checker, exists := ownerChecker[name]
		if !exists {
			return fmt.Errorf("%w: \"%s\"", ErrUnexpectObjectType, name)
		}

		return checker(viewer, obj)
	}
}

func newRoleCheckers(r ports.Role, ownerChecker map[string]OwnerChecker) map[models.RoleEnum]RoleChecker {
	return map[models.RoleEnum]RoleChecker{
		models.RoleEnumUser:    userRoleChecker,
		models.RoleEnumOwner:   ownerRoleChecker(ownerChecker),
		models.RoleEnumManager: roleTypeRoleChecker(r, models.RoleTypeManager),
		models.RoleEnumAdmin:   roleTypeRoleChecker(r, models.RoleTypeAdmin),
	}
}
