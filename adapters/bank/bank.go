package bank

import (
	"auction-back/models"
	"auction-back/ports"
	"errors"
	"fmt"
)

type BankAdapter struct {
	db ports.DB
}

func New(db ports.DB) BankAdapter {
	return BankAdapter{db: db}
}

func (b *BankAdapter) createAccounts(userID string) error {
	nominalAccounts, err := b.db.NominalAccount().Find(ports.NominalAccountFindConfig{})
	if err != nil {
		return fmt.Errorf("db.NominalAccount().Find: %w", err)
	}

	for _, nominalAccount := range nominalAccounts {
		// TODO: call bank service for account creation
		account := models.Account{
			UserID:           userID,
			NominalAccountID: nominalAccount.ID,
		}

		if err := b.db.Account().Create(&account); err != nil {
			return fmt.Errorf("db account create: %w", err)
		}
	}

	return nil
}

func (b *BankAdapter) UserFormApproved(form models.UserForm) error {
	// TODO: send request to bank for create account or update account and
	// if account is succefully created create account in our database and inform client
	// if account cration is rejected decline user form and inform managers

	// FIXME: this code should be in bank service
	_, err := b.db.Account().Take(ports.AccountTakeConfig{
		Filter: &models.AccountsFilter{
			UserIDs: []string{form.UserID},
		},
	})

	if err == nil {
		// TODO: update data in bank
		return nil
	}

	if errors.Is(err, ports.ErrRecordNotFound) {
		return b.createAccounts(form.UserID)
	}

	return fmt.Errorf("db account take: %w", err)
}
