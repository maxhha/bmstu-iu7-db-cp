package token

import (
	"auction-back/db"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Interface interface {
	Create(action db.TokenAction, viewer *db.User, data map[string]interface{}) error
	Activate(action db.TokenAction, token_code string, viewer *db.User) (db.Token, error)
}

type TokenService struct {
	db *gorm.DB
}

func New(db *gorm.DB) TokenService {
	return TokenService{db}
}

func (t *TokenService) send(token db.Token) error {
	return nil
}

func (t *TokenService) Create(action db.TokenAction, viewer *db.User, data map[string]interface{}) error {
	token := db.Token{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)),
		Action:    action,
		Data:      data,
		UserID:    viewer.ID,
	}

	if err := t.db.Create(&token).Error; err != nil {
		return err
	}

	if err := t.send(token); err != nil {
		return fmt.Errorf("send: %w", err)
	}

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

	if token.ActivatedAt.Valid {
		return token, fmt.Errorf("activated")
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
