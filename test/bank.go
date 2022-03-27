package test

import (
	"auction-back/db"

	"github.com/stretchr/testify/mock"
)

type BankPort struct {
	mock.Mock
}

func (b *BankPort) UserFormApproved(form db.UserForm) error {
	args := b.Called(form)
	return args.Error(0)
}
