package test

import (
	"auction-back/db"

	"github.com/stretchr/testify/mock"
)

type TokenPort struct {
	mock.Mock
}

func (t *TokenPort) Create(action db.TokenAction, viewer *db.User, data map[string]interface{}) error {
	args := t.Called(action, viewer, data)
	return args.Error(0)
}

func (t *TokenPort) Activate(action db.TokenAction, token_code string, viewer *db.User) (db.Token, error) {
	args := t.Called(action, token_code, viewer)
	return args.Get(0).(db.Token), args.Error(1)
}
