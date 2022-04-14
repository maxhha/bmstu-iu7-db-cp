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

func (d *auctionDB) Get(id string) (models.Auction, error) {
	obj := Auction{}
	if err := d.db.Take(&obj, "id = ?", id).Error; err != nil {
		return models.Auction{}, fmt.Errorf("take: %w", convertError(err))
	}

	return obj.into(), nil
}

func (d *auctionDB) Update(auction *models.Auction) error {
	if auction == nil {
		return ports.ErrAuctionIsNil
	}

	a := Auction{}
	a.copy(auction)

	if err := d.db.Save(&a).Error; err != nil {
		return fmt.Errorf("save: %w", convertError(err))
	}

	return nil
}

func (d *auctionDB) Create(auction *models.Auction) error {
	if auction == nil {
		return fmt.Errorf("auction is nil")
	}
	a := Auction{}
	a.copy(auction)
	if err := d.db.Create(&a).Error; err != nil {
		return fmt.Errorf("create: %w", convertError(err))
	}

	*auction = a.into()
	return nil
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

func (d *auctionDB) Take(filter *models.AuctionsFilter) (models.Auction, error) {
	query := d.filter(d.db.Model(&Auction{}), filter)
	auction := Auction{}
	if err := query.Take(&auction).Error; err != nil {
		return models.Auction{}, fmt.Errorf("take: %w", convertError(err))
	}

	return auction.into(), nil
}

func (d *auctionDB) Pagination(first *int, after *string, filter *models.AuctionsFilter) (models.AuctionsConnection, error) {
	query := d.filter(d.db.Model(&Auction{}), filter)
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
