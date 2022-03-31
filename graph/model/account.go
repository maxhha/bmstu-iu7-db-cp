package model

import (
	"auction-back/db"
	"fmt"
)

type BankAccount struct {
	db.Account
}

func (BankAccount) IsAccount() {}

type UserAccount struct {
	db.Account
}

func (UserAccount) IsAccount() {}

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
