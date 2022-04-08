package database

import (
	"auction-back/models"
	"database/sql"
	"time"
)

type Auction struct {
	ID           string              `gorm:"default:generated();"`
	State        models.AuctionState `gorm:"default:'CREATED';"`
	ProductID    string
	Product      Product
	SellerID     string
	Seller       User
	BuyerID      string
	Buyer        User
	ScheduledFor time.Time
	StartedAt    sql.NullTime
	FinishedAt   sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
