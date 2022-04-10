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
}

func New(DB ports.DB, token ports.Token, bank ports.Bank, role ports.Role) *Resolver {
	r := Resolver{
		DB:         DB,
		Market:     make(map[string]chan *models.Product),
		MarketLock: sync.Mutex{},
		TokenPort:  token,
		BankPort:   bank,
		RolePort:   role,
	}

	return &r
}
