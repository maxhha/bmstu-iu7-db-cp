package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductImage struct {
	gorm.Model
	ID        string `gorm:"type:varchar(16);"`
	Filename  string
	Path      string
	ProductID string `gorm:"type:varchar(16);"`
	Product   Product
}

type OfferState string

const (
	OfferStateCreated               OfferState = "CREATED"
	OfferStateCancelled             OfferState = "CANCELLED"
	OfferStateTransferringMoney     OfferState = "TRANSFERRING_MONEY"
	OfferStateTransferMoneyFailed   OfferState = "TRANSFER_MONEY_FAILED"
	OfferStateTransferringProduct   OfferState = "TRANSFERRING_PRODUCT"
	OfferStateTransferProductFailed OfferState = "TRANSFER_PRODUCT_FAILED"
	OfferStateSucceeded             OfferState = "SUCCEEDED"
	OfferStateReturningMoney        OfferState = "RETURNING_MONEY"
	OfferStateReturnMoneyFailed     OfferState = "RETURN_MONEY_FAILED"
	OfferStateMoneyReturned         OfferState = "MONEY_RETURNED"
)

func (s *OfferState) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("convert to bytes")
	}

	*s = OfferState(bytes)
	return nil
}

func (s OfferState) Value() (driver.Value, error) {
	return string(s), nil
}

type Offer struct {
	gorm.Model
	ID           string     `gorm:"type:varchar(16);"`
	State        OfferState `gorm:"type:offer_state;"`
	DeleteOnSell bool       `gorm:"default:true"`
	ProductID    string     `gorm:"type:varchar(16);"`
	UserID       string     `gorm:"type:varchar(16);"`
	User         User
	Product      Product
}

type TransactionState string

const (
	TransactionStateCreated    TransactionState = "CREATED"
	TransactionStateCancelled  TransactionState = "CANCELLED"
	TransactionStateProcessing TransactionState = "PROCESSING"
	TransactionStateError      TransactionState = "ERROR"
	TransactionStateSucceeded  TransactionState = "SUCCEEDED"
	TransactionStateFailed     TransactionState = "FAILED"
)

func (s *TransactionState) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("convert to bytes")
	}

	*s = TransactionState(bytes)
	return nil
}

func (s TransactionState) Value() (driver.Value, error) {
	return string(s), nil
}

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
	TransactionTypeBuy        TransactionType = "BUY"
	TransactionTypeFee        TransactionType = "FEE"
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
)

func (t *TransactionType) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("convert to bytes")
	}

	*t = TransactionType(bytes)
	return nil
}

func (t TransactionType) Value() (driver.Value, error) {
	return string(t), nil
}

type TransactionCurrency string

const (
	TransactionCurrencyRub TransactionCurrency = "RUB"
	TransactionCurrencyUsd TransactionCurrency = "USD"
	TransactionCurrencyEur TransactionCurrency = "EUR"
)

func (c *TransactionCurrency) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("convert to bytes")
	}

	*c = TransactionCurrency(bytes)
	return nil
}

func (c TransactionCurrency) Value() (driver.Value, error) {
	return string(c), nil
}

type Transaction struct {
	gorm.Model
	Date          time.Time
	State         TransactionState    `gorm:"type:transaction_state;"`
	Type          TransactionType     `gorm:"type:transaction_type;"`
	Currency      TransactionCurrency `gorm:"type:transaction_currency;"`
	Amount        decimal.Decimal     `gorm:"type:decimal(10,2);"`
	Error         *string
	AccountFromID string `gorm:"type:varchar(16);"`
	AccountToID   string `gorm:"type:varchar(16);"`
	AccountFrom   Account
	AccountTo     Account
	OfferID       string `gorm:"type:varchar(16);"`
	Offer         Offer
}
