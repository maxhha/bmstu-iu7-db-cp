package bank

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"errors"
	"fmt"
)

type BankPort struct {
	db ports.DB
}

func New(db ports.DB) BankPort {
	return BankPort{db}
}

func (b *BankPort) createAccount(userID string) error {
	bank, err := b.db.Bank().Take(ports.BankTakeConfig{
		Names: []string{"fake"},
	})
	if err != nil {
		return fmt.Errorf("db take: %w", err)
	}

	account := models.Account{
		Type:   models.AccountTypeUser,
		UserID: userID,
		BankID: bank.ID,
	}

	if err := b.db.Account().Create(&account); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	return nil
}

func (b *BankPort) UserFormApproved(form models.UserForm) error {
	// TODO: send request to bank for create account or update account and
	// if account is succefully created create account in our database and inform client
	// if account cration is rejected decline user form and inform managers

	// IDEA: create service for each bank. Using config from env or from UserForm
	// select bank services to call.

	// FIXME: this code should be in bank service
	_, err := b.db.Account().Take(ports.AccountTakeConfig{
		UserIDs: []string{form.UserID},
	})

	if err == nil {
		// TODO: update data in bank
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return b.createAccount(form.UserID)
	}

	return fmt.Errorf("db take: %w", err)
}
