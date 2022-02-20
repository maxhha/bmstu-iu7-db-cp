package graph

import (
	"auction-back/graph/model"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Market     map[string]chan *model.Product
	MarketLock sync.Mutex
}

func New() *Resolver {
	r := Resolver{
		Market:     make(map[string]chan *model.Product),
		MarketLock: sync.Mutex{},
	}

	return &r
}
