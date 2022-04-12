package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

type tokenDB struct{ *Database }

func (d *Database) Token() ports.TokenDB { return &tokenDB{d} }

type Token struct {
	ID          uint
	UserID      string
	CreatedAt   time.Time `gorm:"default:now();"`
	ActivatedAt sql.NullTime
	ExpiresAt   time.Time
	Action      models.TokenAction
	Data        datatypes.JSONMap
}

func (t *Token) into() models.Token {
	return models.Token{
		ID:          t.ID,
		UserID:      t.UserID,
		CreatedAt:   t.CreatedAt,
		ActivatedAt: t.ActivatedAt,
		ExpiresAt:   t.ExpiresAt,
		Action:      t.Action,
		Data:        t.Data,
	}
}

func (t *Token) copy(token *models.Token) {
	if token == nil {
		return
	}
	t.ID = token.ID
	t.UserID = token.UserID
	t.CreatedAt = token.CreatedAt
	t.ActivatedAt = token.ActivatedAt
	t.ExpiresAt = token.ExpiresAt
	t.Action = token.Action
	t.Data = token.Data
}

func (d *tokenDB) Create(token *models.Token) error {
	if token == nil {
		return fmt.Errorf("token is nil")
	}
	t := Token{}
	t.copy(token)
	if err := d.db.Create(&t).Error; err != nil {
		return fmt.Errorf("create: %w", convertError(err))
	}

	*token = t.into()
	return nil
}

func (d *tokenDB) Take(config ports.TokenTakeConfig) (models.Token, error) {
	query := &d.db

	if len(config.IDs) > 0 {
		query = query.Where("id IN ?", config.IDs)
	}

	if len(config.UserIDs) > 0 {
		query = query.Where("user_id IN ?", config.UserIDs)
	}

	token := Token{}
	if err := query.Take(&token).Error; err != nil {
		return models.Token{}, fmt.Errorf("take: %w", convertError(err))
	}

	return token.into(), nil
}

func (d *tokenDB) Update(token *models.Token) error {
	if token == nil {
		return fmt.Errorf("token is nil")
	}

	f := Token{}
	f.copy(token)
	if err := d.db.Save(&f).Error; err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (d *tokenDB) GetUser(token models.Token) (models.User, error) {
	return d.User().Get(token.UserID)
}
