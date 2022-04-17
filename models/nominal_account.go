package models

import (
	"database/sql"
	"time"
)

type NominalAccount struct {
	ID            string
	Name          string
	Receiver      string
	AccountNumber string
	BankID        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}
