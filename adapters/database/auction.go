package database

import (
	"auction-back/models"
	"auction-back/ports"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate go run ../../codegen/gormdbops/main.go --out auction_gen.go --model Auction --methods Get,Update,Create,Take,Pagination

type Auction struct {
	ID                string              `gorm:"default:generated();"`
	State             models.AuctionState `gorm:"default:'CREATED';"`
	ProductID         string
	SellerID          string
	SellerAccountID   *string
	BuyerID           *string
	MinAmount         *decimal.Decimal
	Currency          models.CurrencyEnum
	ScheduledStartAt  *time.Time
	ScheduledFinishAt *time.Time
	StartedAt         *time.Time
	FinishedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (a *Auction) into() models.Auction {
	return models.Auction{
		ID:                a.ID,
		State:             a.State,
		ProductID:         a.ProductID,
		SellerID:          a.SellerID,
		SellerAccountID:   a.SellerAccountID,
		BuyerID:           a.BuyerID,
		Currency:          a.Currency,
		MinAmount:         a.MinAmount,
		ScheduledStartAt:  a.ScheduledStartAt,
		ScheduledFinishAt: a.ScheduledFinishAt,
		StartedAt:         a.StartedAt,
		FinishedAt:        a.FinishedAt,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

func (a *Auction) copy(auction *models.Auction) {
	if auction == nil {
		return
	}

	a.ID = auction.ID
	a.State = auction.State
	a.ProductID = auction.ProductID
	a.SellerID = auction.SellerID
	a.SellerAccountID = auction.SellerAccountID
	a.BuyerID = auction.BuyerID
	a.Currency = auction.Currency
	a.ScheduledStartAt = auction.ScheduledStartAt
	a.ScheduledFinishAt = auction.ScheduledFinishAt
	a.StartedAt = auction.StartedAt
	a.FinishedAt = auction.FinishedAt
	a.CreatedAt = auction.CreatedAt
	a.UpdatedAt = auction.UpdatedAt
}

func (d *auctionDB) filter(query *gorm.DB, config *models.AuctionsFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.IDs) > 0 {
		query = query.Where("id IN ?", config.IDs)
	}

	if len(config.SellerIDs) > 0 {
		query = query.Where("seller_id IN ?", config.SellerIDs)
	}

	if len(config.BuyerIDs) > 0 {
		query = query.Where("buyer_id IN ?", config.BuyerIDs)
	}

	if len(config.ProductIDs) > 0 {
		query = query.Where("product_id IN ?", config.ProductIDs)
	}

	if len(config.States) > 0 {
		query = query.Where("state IN ?", config.States)
	}

	return query
}

func (d *auctionDB) LockShare(auction *models.Auction) error {
	if auction == nil {
		return ports.ErrAuctionIsNil
	}
	obj := Auction{}
	err := d.db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).
		Take(&obj, "id = ?", auction.ID).
		Error
	if err != nil {
		return convertError(err)
	}

	*auction = obj.into()
	return nil
}
