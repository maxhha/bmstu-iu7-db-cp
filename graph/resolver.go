package graph

import (
	"auction-back/models"
	"auction-back/ports"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB         ports.DB
	Market     map[string]chan *models.Product
	MarketLock sync.Mutex
	TokenPort  ports.Token
	BankPort   ports.Bank
	RolePort   ports.Role
	DealerPort ports.Dealer
}

func New(DB ports.DB, token ports.Token, bank ports.Bank, role ports.Role, dealer ports.Dealer) *Resolver {
	r := Resolver{
		DB:         DB,
		Market:     make(map[string]chan *models.Product),
		MarketLock: sync.Mutex{},
		TokenPort:  token,
		BankPort:   bank,
		RolePort:   role,
		DealerPort: dealer,
	}

	return &r
}

func (r *Resolver) Tx(fn func(tx ports.DB) error) error {
	tx := r.DB.Tx()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err := fn(tx.DB())

	if err == nil {
		return tx.Commit()
	} else {
		tx.Rollback()
		return err
	}
}
