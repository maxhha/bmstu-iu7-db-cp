package model

import "auction-back/db"

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsOnMarket  bool    `json:"isOnMarket"`
	DB          *db.Product
}

func (p *Product) From(product *db.Product) (*Product, error) {
	p.ID = product.ID
	// p.Name = product.Name
	// p.Description = product.Description
	p.IsOnMarket = product.IsOnMarket
	p.DB = product

	return p, nil
}
