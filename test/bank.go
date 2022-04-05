package test

import (
	"auction-back/models"

	"github.com/stretchr/testify/mock"
)

type BankPort struct {
	mock.Mock
}

func (b *BankPort) UserFormApproved(form models.UserForm) error {
	args := b.Called(form)
	return args.Error(0)
}
