package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type bankDB struct{ *Database }

func (d *Database) Bank() ports.BankDB { return &bankDB{d} }

type Bank struct {
	ID        string `gorm:"default:generated();"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (b *Bank) into() models.Bank {
	return models.Bank{
		ID:        b.ID,
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		DeletedAt: sql.NullTime(b.DeletedAt),
	}
}

func (d *bankDB) Get(id string) (models.Bank, error) {
	obj := Bank{}
	if err := d.db.Take(&obj, "id = ?", id).Error; err != nil {
		return models.Bank{}, err
	}

	return obj.into(), nil
}

func (d *Database) GetAccount(bank models.Bank) (models.BankAccount, error) {
	obj := Account{}
	err := d.db.Take(
		&obj,
		"type = ? AND bank_id = ?",
		models.AccountTypeBank,
		bank.ID).
		Error
	if err != nil {
		return models.BankAccount{}, fmt.Errorf("take: %w", convertError(err))
	}

	account := obj.into()
	conAccount, err := account.ConcreteType()
	if err != nil {
		return models.BankAccount{}, fmt.Errorf("%v concrete type: %w", account, err)
	}

	switch account := conAccount.(type) {
	case models.BankAccount:
		return account, nil
	default:
		return models.BankAccount{}, fmt.Errorf("unexpected bank account type: %s", obj.Type)
	}
}

func (d *bankDB) Take(config ports.BankTakeConfig) (models.Bank, error) {
	query := d.db

	if len(config.IDs) > 0 {
		query = query.Where("id IN ?", config.IDs)
	}

	if len(config.Names) > 0 {
		query = query.Where("name IN ?", config.Names)
	}

	bank := Bank{}
	if err := query.Take(&bank).Error; err != nil {
		return models.Bank{}, fmt.Errorf("take: %w", convertError(err))
	}

	return bank.into(), nil
}
