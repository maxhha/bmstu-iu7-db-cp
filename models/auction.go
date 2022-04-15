package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Auction struct {
	ID                string `json:"id"`
	State             AuctionState
	ProductID         string
	SellerID          string
	BuyerID           *string
	MinAmount         *decimal.Decimal
	Currency          CurrencyEnum `json:"currency"`
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
