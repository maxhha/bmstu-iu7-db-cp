package test

import (
	"auction-back/models"

	"github.com/stretchr/testify/mock"
)

type RolePort struct {
	mock.Mock
}

func (r *RolePort) HasRole(roleType models.RoleType, viewer models.User) error {
	args := r.Called(roleType, viewer)
	return args.Error(0)
}
