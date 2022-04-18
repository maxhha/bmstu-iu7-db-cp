package database

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate go run ../../codegen/gormdbops/main.go --out offer_gen.go --model Offer --methods Get,Update,Create,Pagination

type Offer struct {
	ID        string            `gorm:"default:generated();"`
	State     models.OfferState `gorm:"default:'CREATED';"`
	AuctionID string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var offerFieldToColumn = map[ports.OfferField]string{
	ports.OfferFieldCreatedAt: "created_at",
}

func (o *Offer) into() models.Offer {
	return models.Offer{
		ID:        o.ID,
		State:     o.State,
		AuctionID: o.AuctionID,
		UserID:    o.UserID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func (o *Offer) copy(offer *models.Offer) {
	if offer == nil {
		return
	}

	o.ID = offer.ID
	o.State = offer.State
	o.AuctionID = offer.AuctionID
	o.UserID = offer.UserID
	o.CreatedAt = offer.CreatedAt
	o.UpdatedAt = offer.UpdatedAt
}

func (d *offerDB) filter(query *gorm.DB, config *models.OffersFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.IDs) > 0 {
		query = query.Where("id IN ?", config.IDs)
	}

	if len(config.AuctionIDs) > 0 {
		query = query.Where("auction_id IN ?", config.AuctionIDs)
	}

	if len(config.States) > 0 {
		query = query.Where("state IN ?", config.States)
	}

	if len(config.UserIDs) > 0 {
		query = query.Where("user_id IN ?", config.UserIDs)
	}

	return query
}

func (d *offerDB) Take(config ports.OfferTakeConfig) (models.Offer, error) {
	query := d.filter(d.db, config.Filter)

	if config.OrderBy != "" {
		column, ok := offerFieldToColumn[config.OrderBy]
		if !ok {
			return models.Offer{}, fmt.Errorf("unknown field '%s'", config.OrderBy)
		}

		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   config.OrderDesc,
		})
	}

	offer := Offer{}
	if err := query.Take(&offer).Error; err != nil {
		return models.Offer{}, fmt.Errorf("take: %w", convertError(err))
	}

	return offer.into(), nil
}
