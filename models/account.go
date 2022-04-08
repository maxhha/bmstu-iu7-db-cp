package models

import (
	"database/sql"
	"fmt"
	"time"
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
	ID        string `json:"id"`
	Type      AccountType
	UserID    string
	BankID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

func (a *Account) String() string {
	return fmt.Sprintf("Account[id=%s]", a.ID)
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
