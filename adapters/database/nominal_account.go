package database

import (
	"auction-back/models"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

//go:generate go run ../../codegen/gormdbops/main.go --out nominal_account_gen.go --model NominalAccount --methods Get,Take,Create,Update,Pagination

type NominalAccount struct {
	ID            string `gorm:"default:generated();"`
	Name          string
	Receiver      string
	AccountNumber string
	BankID        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (a *NominalAccount) into() models.NominalAccount {
	return models.NominalAccount{
		ID:            a.ID,
		Name:          a.Name,
		Receiver:      a.Receiver,
		AccountNumber: a.AccountNumber,
		BankID:        a.BankID,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
		DeletedAt:     sql.NullTime(a.DeletedAt),
	}
}

func (a *NominalAccount) copy(account *models.NominalAccount) {
	if account == nil {
		return
	}

	a.ID = account.ID
	a.Name = account.Name
	a.Receiver = account.Receiver
	a.AccountNumber = account.AccountNumber
	a.BankID = account.BankID
	a.CreatedAt = account.CreatedAt
	a.UpdatedAt = account.UpdatedAt
	a.DeletedAt = gorm.DeletedAt(account.DeletedAt)
}

func (d *nominalAccountDB) filter(query *gorm.DB, config *models.NominalAccountsFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.IDs) > 0 {
		query = query.Where("id IN ?", config.IDs)
	}

	if len(config.BankIDs) > 0 {
		query = query.Where("bank_id IN ?", config.BankIDs)
	}

	if config.Name != nil && *config.Name != "" {
		query = query.Where("name ~ ?", config.Name)
	}

	return query
}
