package models

import (
	"time"
)

type Auction struct {
	ID                string `json:"id"`
	State             AuctionState
	ProductID         string
	SellerID          string
	BuyerID           *string
	MinMoney          *Money `json:"minMoney"`
	ScheduledStartAt  *time.Time
	ScheduledFinishAt *time.Time
	StartedAt         *time.Time
	FinishedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (a *Auction) IsEditable() bool {
	return a.State == AuctionStateCreated
}
