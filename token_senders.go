package main

import (
	"auction-back/adapters/token_sender"
	"auction-back/db"
	"fmt"
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

func phoneTokenSender() *token_sender.TokenSender {
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
	}

	config.DataGetters = map[db.TokenAction]token_sender.DataGetter{
		db.TokenActionSetUserEmail: dataWithTokenId,
	}

	return token_sender.New(config)
}
