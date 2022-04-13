package database

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type auctionDB struct{ *Database }

func (d *Database) Auction() ports.AuctionDB { return &auctionDB{d} }

type Auction struct {
	ID                string              `gorm:"default:generated();"`
	State             models.AuctionState `gorm:"default:'CREATED';"`
	ProductID         string
	SellerID          string
	BuyerID           *string
	minAmount         *decimal.Decimal
	minAmountCurrency *models.CurrencyEnum
	ScheduledStartAt  *time.Time
	ScheduledFinishAt *time.Time
	StartedAt         *time.Time
	FinishedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (a *Auction) into() models.Auction {
	var money *models.Money

	if a.minAmount != nil && a.minAmountCurrency != nil {
		money = &models.Money{
			Amount:   *a.minAmount,
			Currency: *a.minAmountCurrency,
		}
	}

	return models.Auction{
		ID:                a.ID,
		State:             a.State,
		ProductID:         a.ProductID,
		SellerID:          a.SellerID,
		BuyerID:           a.BuyerID,
		MinMoney:          money,
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
	a.BuyerID = auction.BuyerID

	if auction.MinMoney != nil {
		a.minAmount = &auction.MinMoney.Amount
		a.minAmountCurrency = &auction.MinMoney.Currency
	}

	a.ScheduledStartAt = auction.ScheduledStartAt
	a.ScheduledFinishAt = auction.ScheduledFinishAt
	a.StartedAt = auction.StartedAt
	a.FinishedAt = auction.FinishedAt
	a.CreatedAt = auction.CreatedAt
	a.UpdatedAt = auction.UpdatedAt
}

func (d *auctionDB) Create(auction *models.Auction) error {
	if auction == nil {
		return fmt.Errorf("auction is nil")
	}
	a := Auction{}
	a.copy(auction)
	if err := d.db.Create(&a).Error; err != nil {
		return fmt.Errorf("create: %w", err)
	}

	*auction = a.into()
	return nil
}

func (d *auctionDB) filter(query *gorm.DB, config *models.AuctionsFilter) *gorm.DB {
	if config == nil {
		return query
	}

	return query
}

func (d *auctionDB) Pagination(first *int, after *string, filter *models.AuctionsFilter) (models.AuctionsConnection, error) {
	query := d.filter(d.db.Model(&Product{}), filter)
	query, err := paginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return models.AuctionsConnection{}, fmt.Errorf("pagination: %w", err)
	}

	var auctions []Auction
	if err := query.Find(&auctions).Error; err != nil {
		return models.AuctionsConnection{}, fmt.Errorf("find: %w", convertError(err))
	}

	if len(auctions) == 0 {
		return models.AuctionsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.AuctionsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(auctions) > *first
		auctions = auctions[:len(auctions)-1]
	}

	edges := make([]*models.AuctionsConnectionEdge, 0, len(auctions))

	for _, obj := range auctions {
		node := obj.into()
		edges = append(edges, &models.AuctionsConnectionEdge{
			Cursor: node.ID,
			Node:   &node,
		})
	}

	return models.AuctionsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &auctions[0].ID,
			EndCursor:   &auctions[len(auctions)-1].ID,
		},
		Edges: edges,
	}, nil
}
