package models

import "time"

type Offer struct {
	ID        string
	State     OfferState
	AuctionID string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
