package graph

import (
	"auction-back/graph/model"
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
}

func New(db *gorm.DB) *Resolver {
	r := Resolver{
		DB:         db,
		Market:     make(map[string]chan *model.Product),
		MarketLock: sync.Mutex{},
	}

	return &r
}
