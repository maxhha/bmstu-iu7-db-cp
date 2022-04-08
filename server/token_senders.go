package server

import (
	"auction-back/adapters/token_sender"
	"auction-back/models"
	"auction-back/ports"
	"fmt"
)

func dataWithTokenId(token models.Token) (map[string]string, error) {
	return map[string]string{
		"token": fmt.Sprintf("%06d", token.ID),
	}, nil
}

func emailTokenSender() *token_sender.TokenSender {
	config := token_sender.Config{
		Name:              "email",
		AddressEnvVarName: "EMAIL_NOTIFIER_ADDRESS",
	}

	config.ReceiverGetters = map[models.TokenAction]token_sender.ReceiverGetter{
		models.TokenActionSetUserEmail: func(token models.Token) (string, error) {
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

	config.DataGetters = map[models.TokenAction]token_sender.DataGetter{
		models.TokenActionSetUserEmail: dataWithTokenId,
	}

	return token_sender.New(config)
}

func phoneTokenSender(DB ports.DB) *token_sender.TokenSender {
	config := token_sender.Config{
		Name:              "phone",
		AddressEnvVarName: "PHONE_NOTIFIER_ADDRESS",
	}

	config.ReceiverGetters = map[models.TokenAction]token_sender.ReceiverGetter{
		models.TokenActionSetUserPhone: func(token models.Token) (string, error) {
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
		models.TokenActionModerateUserForm: func(token models.Token) (string, error) {
			user, err := DB.Token().GetUser(token)
			if err != nil {
				return "", fmt.Errorf("db token get user: %w", err)
			}

			form, err := DB.User().MostRelevantUserForm(user)

			if err != nil {
				return "", fmt.Errorf("last relevant user form: %w", err)
			}

			if form.Phone == nil {
				return "", fmt.Errorf("user form phone is nil")
			}

			return *form.Phone, nil
		},
	}

	config.DataGetters = map[models.TokenAction]token_sender.DataGetter{
		models.TokenActionSetUserPhone:     dataWithTokenId,
		models.TokenActionModerateUserForm: dataWithTokenId,
	}

	return token_sender.New(config)
}
