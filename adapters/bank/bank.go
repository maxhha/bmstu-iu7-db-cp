package bank

import (
	"auction-back/models"
	"auction-back/ports"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
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
		_, err := b.CreateAccount(userID, nominalAccount.ID)
		if err != nil {
			return fmt.Errorf("b.CreateAccount: %w", err)
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

func (b *BankAdapter) CreateAccount(userID string, nominalAccountID string) (models.Account, error) {
	// TODO: call bank service for account creation
	account := models.Account{
		UserID:           userID,
		NominalAccountID: nominalAccountID,
	}

	if err := b.db.Account().Create(&account); err != nil {
		return account, fmt.Errorf("db account create: %w", err)
	}

	return account, nil
}

func (b *BankAdapter) ProcessTransactions(transacions []models.Transaction) error {
	var errs error

	for _, transaction := range transacions {
		if transaction.State == models.TransactionStateCreated || transaction.State == models.TransactionStateError {
			transaction.State = models.TransactionStateProcessing
			if err := b.db.Transaction().Update(&transaction); err != nil {
				errs = multierror.Append(errs, fmt.Errorf("b.db.Transaction().Update(id=%d): %w", transaction.ID, err))
			}
		}
	}

	return errs
}
