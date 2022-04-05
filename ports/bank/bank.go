package bank

import (
	"auction-back/models"
	"fmt"

	"gorm.io/gorm"
)

type Interface interface {
	UserFormApproved(form models.UserForm) error
}

type BankPort struct {
	db *gorm.DB
}

func New(db *gorm.DB) BankPort {
	return BankPort{db}
}

func (b *BankPort) createAccount(userID string) error {
	bank := models.Bank{}

	if err := b.db.Take(&bank, "name = 'fake'").Error; err != nil {
		return fmt.Errorf("take bank: %w", err)
	}

	account := models.Account{
		Type:   models.AccountTypeUser,
		UserID: userID,
		BankID: bank.ID,
	}

	if err := b.db.Create(&account).Error; err != nil {
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
	account := models.Account{}

	err := b.db.Take(&account, "user_id = ?", form.UserID).Error

	if err == nil {
		// TODO: update data in bank
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		return b.createAccount(form.UserID)
	}

	return fmt.Errorf("take: %w", err)
}
