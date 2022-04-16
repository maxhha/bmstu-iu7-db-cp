package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID            int
	Date          *time.Time
	State         TransactionState
	Type          TransactionType
	Currency      CurrencyEnum
	Amount        decimal.Decimal
	Error         *string
	AccountFromID string
	AccountToID   string
	OfferID       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
