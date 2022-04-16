package token_mock

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"
)

type TokenPort struct {
	db ports.DB
}

func New(db ports.DB) TokenPort {
	return TokenPort{db}
}

func (t *TokenPort) Create(action models.TokenAction, viewer models.User, data map[string]interface{}) error {
	token := models.Token{
		ExpiresAt: time.Now().UTC().Add(time.Hour * time.Duration(1)),
		Action:    action,
		Data:      data,
		UserID:    viewer.ID,
	}

	if err := t.db.Token().Create(&token); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	return nil
}

func (t *TokenPort) Activate(action models.TokenAction, token_code string, viewer models.User) (models.Token, error) {
	token, err := t.db.Token().Take(ports.TokenTakeConfig{
		UserIDs:   []string{viewer.ID},
		Actions:   []models.TokenAction{action},
		OrderBy:   ports.TokenFieldCreatedAt,
		OrderDesc: true,
	})
	if err != nil {
		return token, fmt.Errorf("db take: %w", err)
	}

	if token.ActivatedAt.Valid {
		return token, fmt.Errorf("already activated")
	}

	token.ActivatedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	if err := t.db.Token().Update(&token); err != nil {
		return token, fmt.Errorf("save: %w", err)
	}

	return token, nil
}
