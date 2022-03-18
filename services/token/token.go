package token

import (
	"auction-back/db"
	"fmt"
)

type Interface interface {
	Validate(action db.TokenAction, data map[string]interface{}) error
	Send(token db.Token) error
}

type TokenService struct {
}

func New() TokenService {
	return TokenService{}
}

var validateTokenData = map[db.TokenAction]func(data map[string]interface{}) error{
	db.TokenActionApproveUserEmail: func(data map[string]interface{}) error {
		email, found := data["email"]

		if !found {
			return fmt.Errorf("no email in data")
		}

		_, ok := email.(string)
		if !ok {
			return fmt.Errorf("email in data is not string")
		}

		return nil
	},
	db.TokenActionApproveUserPhone: func(data map[string]interface{}) error {
		phone, found := data["phone"]

		if !found {
			return fmt.Errorf("no phone in data")
		}

		_, ok := phone.(string)
		if !ok {
			return fmt.Errorf("phone in data is not string")
		}

		return nil
	},
}

func (t *TokenService) Validate(action db.TokenAction, data map[string]interface{}) error {
	validate, found := validateTokenData[action]
	if !found {
		return fmt.Errorf("not found validator for action")
	}

	if err := validate(data); err != nil {
		return err
	}

	return nil
}

func (t *TokenService) Send(token db.Token) error {
	return nil
}
