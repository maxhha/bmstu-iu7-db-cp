package token

import (
	"auction-back/db"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Interface interface {
	Validate(action db.TokenAction, data map[string]interface{}) error
	Send(token db.Token) error
	Activate(action db.TokenAction, token_code string, viewer *db.User) (db.Token, error)
}

type TokenService struct {
	db *gorm.DB
}

func New(db *gorm.DB) TokenService {
	return TokenService{db}
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

func (t *TokenService) Activate(action db.TokenAction, token_code string, viewer *db.User) (db.Token, error) {
	if viewer == nil {
		return db.Token{}, fmt.Errorf("unauthorized")
	}

	token := db.Token{}

	if err := t.db.Take(&token, "id = ? and user_id = ?", token_code, viewer.ID).Error; err != nil {
		return db.Token{}, fmt.Errorf("take: %w", err)
	}

	if token.Action != action {
		return db.Token{}, fmt.Errorf("action not match")
	}

	token.ActivatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err := t.db.Save(&token).Error; err != nil {
		return db.Token{}, fmt.Errorf("save: %w", err)
	}

	return token, nil
}
