package model

import (
	"auction-back/db"
	"fmt"
)

type BankAccount struct {
	db.Account
}

func (BankAccount) IsAccount() {}

func (a *BankAccount) AccountPtr() *db.Account {
	if a == nil {
		return nil
	}

	return &a.Account
}

type UserAccount struct {
	db.Account
}

func (UserAccount) IsAccount() {}

func (a *UserAccount) AccountPtr() *db.Account {
	if a == nil {
		return nil
	}

	return &a.Account
}

func AccountFromDBAccount(obj db.Account) (Account, error) {
	switch {
	case obj.Type == db.AccountTypeBank:
		return BankAccount{Account: obj}, nil
	case obj.Type == db.AccountTypeUser:
		return UserAccount{Account: obj}, nil
	default:
		return nil, fmt.Errorf("unexpected account type %s", obj.Type)
	}
}
