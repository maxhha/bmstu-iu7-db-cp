package server

import (
	"auction-back/adapters/token_sender"
	"auction-back/db"
	"fmt"

	"gorm.io/gorm"
)

func dataWithTokenId(token db.Token) (map[string]string, error) {
	return map[string]string{
		"token": fmt.Sprintf("%06d", token.ID),
	}, nil
}

func emailTokenSender() *token_sender.TokenSender {
	config := token_sender.Config{
		Name:              "email",
		AddressEnvVarName: "EMAIL_NOTIFIER_ADDRESS",
	}

	config.ReceiverGetters = map[db.TokenAction]token_sender.ReceiverGetter{
		db.TokenActionSetUserEmail: func(token db.Token) (string, error) {
			email, ok := token.Data["email"]
			if !ok {
				return "", fmt.Errorf("no email in token data")
			}

			str, ok := email.(string)
			if !ok {
				return "", fmt.Errorf("fail convert to string of %v", email)
			}

			return str, nil
		},
	}

	config.DataGetters = map[db.TokenAction]token_sender.DataGetter{
		db.TokenActionSetUserEmail: dataWithTokenId,
	}

	return token_sender.New(config)
}

func phoneTokenSender(DB *gorm.DB) *token_sender.TokenSender {
	config := token_sender.Config{
		Name:              "phone",
		AddressEnvVarName: "PHONE_NOTIFIER_ADDRESS",
	}

	config.ReceiverGetters = map[db.TokenAction]token_sender.ReceiverGetter{
		db.TokenActionSetUserPhone: func(token db.Token) (string, error) {
			phone, ok := token.Data["phone"]
			if !ok {
				return "", fmt.Errorf("no phone in token data")
			}

			str, ok := phone.(string)
			if !ok {
				return "", fmt.Errorf("fail convert to string of %v", phone)
			}

			return str, nil
		},
		db.TokenActionModerateUserForm: func(token db.Token) (string, error) {
			if err := token.EnsureFillUser(DB); err != nil {
				return "", err
			}

			form, err := token.User.MostRelevantUserForm(DB)

			if err != nil {
				return "", fmt.Errorf("last relevant user form: %w", err)
			}

			if form.Phone == nil {
				return "", fmt.Errorf("user form phone is nil")
			}

			return *form.Phone, nil
		},
	}

	config.DataGetters = map[db.TokenAction]token_sender.DataGetter{
		db.TokenActionSetUserPhone:     dataWithTokenId,
		db.TokenActionModerateUserForm: dataWithTokenId,
	}

	return token_sender.New(config)
}
