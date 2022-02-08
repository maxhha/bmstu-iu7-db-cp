package model

import "auction-back/db"

type Offer struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"createdAt"`
	DB        *db.Offer
}

func (o *Offer) From(offer *db.Offer) (*Offer, error) {
	o.ID = offer.ID
	o.Amount = offer.Amount
	o.CreatedAt = offer.CreatedAt.String()
	o.DB = offer

	return o, nil
}
