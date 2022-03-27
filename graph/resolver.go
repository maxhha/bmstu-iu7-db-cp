package graph

import (
	"auction-back/graph/model"
	"auction-back/ports/bank"
	"auction-back/ports/token"
	"sync"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB         *gorm.DB
	Market     map[string]chan *model.Product
	MarketLock sync.Mutex
	Token      token.Interface
	Bank       bank.Interface
}

func New(db *gorm.DB, token token.Interface, bank bank.Interface) *Resolver {
	r := Resolver{
		DB:         db,
		Market:     make(map[string]chan *model.Product),
		MarketLock: sync.Mutex{},
		Token:      token,
		Bank:       bank,
	}

	return &r
}
