package server

import (
	"auction-back/adapters/token_sender"
	"auction-back/models"
	"auction-back/ports"
	"errors"
	"fmt"
	"strings"
)

var ErrUserEmailIsNil = errors.New("user form email is nil")
var ErrUserPhoneIsNil = errors.New("user form phone is nil")
var ErrTokenDataNoField = errors.New("token data not have field")

func userEmailReceiverGetter(DB ports.DB) token_sender.ReceiverGetter {
	return func(token models.Token) (string, error) {
		user, err := DB.Token().GetUser(token)
		if err != nil {
			return "", fmt.Errorf("db token get user: %w", err)
		}

		form, err := DB.User().MostRelevantUserForm(user)
		if err != nil {
			return "", fmt.Errorf("last relevant user form: %w", err)
		}

		if form.Email == nil {
			return "", ErrUserEmailIsNil
		}

		return *form.Email, nil
	}
}

func tokenEmailReceiverGetter(getUserEmail token_sender.ReceiverGetter) token_sender.ReceiverGetter {
	return func(token models.Token) (string, error) {
		userEmail, err := getUserEmail(token)
		if err == nil {
			return userEmail, nil
		} else if !errors.Is(err, ErrUserEmailIsNil) &&
			!strings.Contains(err.Error(), "last relevant user form: take: not found") {
			return "", err
		}

		email, ok := token.Data["email"]
		if !ok {
			return "", fmt.Errorf("%w: email", ErrTokenDataNoField)
		}

		str, ok := email.(string)
		if !ok {
			return "", fmt.Errorf("fail convert to string of %v", email)
		}

		return str, nil
	}
}

func userPhoneReceiverGetter(DB ports.DB) token_sender.ReceiverGetter {
	return func(token models.Token) (string, error) {
		user, err := DB.Token().GetUser(token)
		if err != nil {
			return "", fmt.Errorf("db token get user: %w", err)
		}

		form, err := DB.User().MostRelevantUserForm(user)

		if err != nil {
			return "", fmt.Errorf("last relevant user form: %w", err)
		}

		if form.Phone == nil {
			return "", ErrUserPhoneIsNil
		}

		return *form.Phone, nil
	}
}

func tokenPhoneReceiverGetter(getUserPhone token_sender.ReceiverGetter) token_sender.ReceiverGetter {
	return func(token models.Token) (string, error) {
		userPhone, err := getUserPhone(token)
		if err == nil {
			return userPhone, nil
		} else if !errors.Is(err, ErrUserPhoneIsNil) &&
			!strings.Contains(err.Error(), "last relevant user form: take: not found") {
			return "", err
		}

		phone, ok := token.Data["phone"]
		if !ok {
			return "", fmt.Errorf("%w: phone", ErrTokenDataNoField)
		}

		str, ok := phone.(string)
		if !ok {
			return "", fmt.Errorf("fail convert to string of %v", phone)
		}

		return str, nil
	}
}

func dataWithTokenId(token models.Token) (map[string]string, error) {
	return map[string]string{
		"token": fmt.Sprintf("%06d", token.ID),
	}, nil
}

func emailTokenSender(DB ports.DB) *token_sender.TokenSender {
	config := token_sender.Config{
		Name:              "email",
		AddressEnvVarName: "EMAIL_NOTIFIER_ADDRESS",
	}

	getUserEmail := userEmailReceiverGetter(DB)

	config.ReceiverGetters = map[models.TokenAction]token_sender.ReceiverGetter{
		models.TokenActionSetUserEmail: tokenEmailReceiverGetter(getUserEmail),
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

	getUserPhone := userPhoneReceiverGetter(DB)

	config.ReceiverGetters = map[models.TokenAction]token_sender.ReceiverGetter{
		models.TokenActionSetUserPhone:     tokenPhoneReceiverGetter(getUserPhone),
		models.TokenActionModerateUserForm: getUserPhone,
		models.TokenActionModerateProduct:  getUserPhone,
	}

	config.DataGetters = map[models.TokenAction]token_sender.DataGetter{
		models.TokenActionSetUserPhone:     dataWithTokenId,
		models.TokenActionModerateUserForm: dataWithTokenId,
		models.TokenActionModerateProduct:  dataWithTokenId,
	}

	return token_sender.New(config)
}
