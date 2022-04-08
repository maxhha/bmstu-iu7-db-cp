package models

import (
	"database/sql"
	"time"
)

type AuctionState string

const (
	AuctionStateCreated   AuctionState = "CREATED"
	AuctionStateStarted   AuctionState = "STARTED"
	AuctionStateFinished  AuctionState = "FINISHED"
	AuctionStateFailed    AuctionState = "FAILED"
	AuctionStateSucceeded AuctionState = "SUCCEEDED"
)

type Auction struct {
	ID           string `json:"id"`
	State        AuctionState
	ProductID    string
	SellerID     string
	BuyerID      string
	ScheduledFor time.Time
	StartedAt    sql.NullTime
	FinishedAt   sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
