package test

import (
	"auction-back/db"

	"github.com/stretchr/testify/mock"
)

type RolePort struct {
	mock.Mock
}

func (r *RolePort) HasRole(roleType db.RoleType, viewer *db.User) error {
	args := r.Called(roleType, viewer)
	return args.Error(0)
}
