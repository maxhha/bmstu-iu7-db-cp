package test

import (
	"auction-back/models"

	"github.com/stretchr/testify/mock"
)

type TokenPort struct {
	mock.Mock
}

func (t *TokenPort) Create(action models.TokenAction, viewer models.User, data map[string]interface{}) error {
	args := t.Called(action, viewer, data)
	return args.Error(0)
}

func (t *TokenPort) Activate(action models.TokenAction, token_code string, viewer models.User) (models.Token, error) {
	args := t.Called(action, token_code, viewer)
	return args.Get(0).(models.Token), args.Error(1)
}
