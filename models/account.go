package models

import (
	"fmt"

	"gorm.io/gorm"
)

type AccountInterface interface {
	IsAccount()
}

type AccountType string

const (
	AccountTypeUser AccountType = "USER"
	AccountTypeBank AccountType = "BANK"
)

type Account struct {
	gorm.Model
	ID     string `gorm:"default:generated();" json:"id"`
	Type   AccountType
	UserID string
	User   User
	BankID string
	Bank   Bank
}

func (a *Account) EnsureFillBank(db *gorm.DB) error {
	if a.BankID == a.Bank.ID {
		return nil
	}

	return db.Take(&a.Bank, "id = ?", a.BankID).Error
}

func (obj *Account) ConcreteType() (AccountInterface, error) {
	switch {
	case obj.Type == AccountTypeBank:
		return BankAccount{Account: *obj}, nil
	case obj.Type == AccountTypeUser:
		return UserAccount{Account: *obj}, nil
	default:
		return nil, fmt.Errorf("unexpected account type %s", obj.Type)
	}
}

type BankAccount struct {
	Account
}

func (BankAccount) IsAccount() {}

func (a *BankAccount) AccountPtr() *Account {
	if a == nil {
		return nil
	}

	return &a.Account
}

type UserAccount struct {
	Account
}

func (UserAccount) IsAccount() {}

func (a *UserAccount) AccountPtr() *Account {
	if a == nil {
		return nil
	}

	return &a.Account
}
