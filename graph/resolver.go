package graph

import (
	"auction-back/graph/model"
	"auction-back/services/token"
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
}

func New(db *gorm.DB, token token.Interface) *Resolver {
	r := Resolver{
		DB:         db,
		Market:     make(map[string]chan *model.Product),
		MarketLock: sync.Mutex{},
		Token:      token,
	}

	return &r
}
