package db

import "gorm.io/gorm"

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
