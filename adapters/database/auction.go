package database

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"
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

func (a *Auction) copy(obj *models.Auction) {
	if obj == nil {
		return
	}

	a.ID = obj.ID
	a.State = obj.State
	a.ProductID = obj.ProductID
	a.SellerID = obj.SellerID
	a.SellerAccountID = obj.SellerAccountID
	a.BuyerID = obj.BuyerID
	a.Currency = obj.Currency
	a.MinAmount = obj.MinAmount
	a.ScheduledStartAt = obj.ScheduledStartAt
	a.ScheduledFinishAt = obj.ScheduledFinishAt
	a.StartedAt = obj.StartedAt
	a.FinishedAt = obj.FinishedAt
	a.CreatedAt = obj.CreatedAt
	a.UpdatedAt = obj.UpdatedAt
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

	if config.ScheduledStartAt != nil {
		if config.StartedAt.From != nil {
			query = query.Where("started_at >= ?", config.StartedAt.From)
		}

		if config.StartedAt.To != nil {
			query = query.Where("started_at < ?", config.StartedAt.To)
		}
	}

	if config.ScheduledStartAt != nil {
		if config.ScheduledStartAt.From != nil {
			query = query.Where("scheduled_start_at >= ?", config.ScheduledStartAt.From)
		}

		if config.ScheduledStartAt.To != nil {
			query = query.Where("scheduled_start_at < ?", config.ScheduledStartAt.To)
		}
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

func (d *auctionDB) forFinishAuctionsQuery(tx *gorm.DB, time_gap string, default_duration string) *gorm.DB {
	return tx.Clauses(clause.Locking{
		Strength: "UPDATE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).
		Where("state = ?", models.AuctionStateStarted).
		Where("NOW() - ?::interval > ( SELECT COALESCE(MAX(offers.created_at), auctions.started_at) FROM offers WHERE offers.auction_id = auctions.id )", time_gap).
		Where("NOW() > COALESCE(auctions.scheduled_finish_at, auctions.started_at + ?::interval)", default_duration)
}

func (d *auctionDB) FindAndSetFinish(config ports.FindAndSetFinishConfig) ([]models.Auction, error) {
	var auctions []models.Auction

	err := d.db.Transaction(func(tx *gorm.DB) error {
		var objs []Auction
		query := d.filter(tx, config.Filter)
		query = d.forFinishAuctionsQuery(query, config.TimeGapFromLastOffer, config.DefaultDuration)
		err := query.Find(&objs).Error

		if err != nil {
			return fmt.Errorf("tx.Find: %w", err)
		}

		auctions = make([]models.Auction, 0, len(objs))
		now := time.Now().UTC()

		for _, a := range objs {
			a.State = models.AuctionStateFinished
			a.FinishedAt = &now
			err = tx.Updates(Auction{ID: a.ID, State: a.State, FinishedAt: a.FinishedAt}).Error
			if err != nil {
				return fmt.Errorf("tx.Update(ID=%s): %w", a.ID, err)
			}
			auctions = append(auctions, a.into())
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("d.db.Transaction: %w", err)
	}

	return auctions, nil
}

func (d *auctionDB) GetTopOffer(auction models.Auction) (models.Offer, error) {
	query := d.db.Model(&Offer{}).Where("offers.auction_id = ?", auction.ID)
	query = d.DB().Offer().(*offerDB).offersNumberedQuery(query)

	offer := Offer{}
	if err := d.db.Model(&offer).
		Joins("JOIN ( ? ) ofd ON offers.id = ofd.offer_id AND ofd.offer_n = 1", query).
		Take(&offer).Error; err != nil {
		return models.Offer{}, fmt.Errorf("take top offer: %w", convertError(err))
	}

	return offer.into(), nil
}
