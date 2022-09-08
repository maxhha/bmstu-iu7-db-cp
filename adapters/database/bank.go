package database

import (
	"auction-back/models"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

//go:generate go run ../../codegen/gormdbops/main.go --out bank_gen.go --model Bank --methods Get,Take,Create,Update,Pagination

type Bank struct {
	ID                   string `gorm:"default:generated();"`
	Name                 string
	Bic                  string
	CorrespondentAccount string
	Inn                  string
	Kpp                  string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
}

func (b *Bank) into() models.Bank {
	return models.Bank{
		ID:                   b.ID,
		Name:                 b.Name,
		Bic:                  b.Bic,
		CorrespondentAccount: b.CorrespondentAccount,
		Inn:                  b.Inn,
		Kpp:                  b.Kpp,
		CreatedAt:            b.CreatedAt,
		UpdatedAt:            b.UpdatedAt,
		DeletedAt:            sql.NullTime(b.DeletedAt),
	}
}

func (b *Bank) copy(bank *models.Bank) {
	if bank == nil {
		return
	}

	b.ID = bank.ID
	b.Name = bank.Name
	b.Bic = bank.Bic
	b.CorrespondentAccount = bank.CorrespondentAccount
	b.Inn = bank.Inn
	b.Kpp = bank.Kpp
	b.CreatedAt = bank.CreatedAt
	b.UpdatedAt = bank.UpdatedAt
	b.DeletedAt = gorm.DeletedAt(bank.DeletedAt)
}

func (d *bankDB) filter(query *gorm.DB, config *models.BanksFilter) *gorm.DB {
	if config == nil {
		return query
	}

	// panic("unimplimented!")
	return query
}
