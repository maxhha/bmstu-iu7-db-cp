package bank

import (
	"auction-back/models"
	"auction-back/ports"
	"errors"
	"fmt"
	"log"
	"os"
)

type BankAdapter struct {
	db                    ports.DB
	defaultNominalAccount string
}

func New(db ports.DB) BankAdapter {
	defaultNominalAccount, exists := os.LookupEnv("BANK_ADAPTER_DEFAULT_NOMINAL_ACCOUNT")
	if !exists {
		log.Fatalln("BANK_ADAPTER_DEFAULT_NOMINAL_ACCOUNT is not set in environment variables")
	}

	return BankAdapter{db: db, defaultNominalAccount: defaultNominalAccount}
}

func (b *BankAdapter) createAccount(userID string) error {
	nominalAccount, err := b.db.NominalAccount().Take(ports.NominalAccountTakeConfig{
		Filter: &models.NominalAccountsFilter{
			Name: &b.defaultNominalAccount,
		},
	})
	if err != nil {
		return fmt.Errorf("db nominal account take: %w", err)
	}

	account := models.Account{
		UserID:           userID,
		NominalAccountID: nominalAccount.ID,
	}

	if err := b.db.Account().Create(&account); err != nil {
		return fmt.Errorf("db account create: %w", err)
	}

	return nil
}

func (b *BankAdapter) UserFormApproved(form models.UserForm) error {
	// TODO: send request to bank for create account or update account and
	// if account is succefully created create account in our database and inform client
	// if account cration is rejected decline user form and inform managers

	// IDEA: create service for each bank. Using config from env or from UserForm
	// select bank services to call.

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
		return b.createAccount(form.UserID)
	}

	return fmt.Errorf("db account take: %w", err)
}
