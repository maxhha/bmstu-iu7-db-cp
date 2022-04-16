package database

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type offerDB struct{ *Database }

func (d *Database) Offer() ports.OfferDB { return &offerDB{d} }

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

func (d *offerDB) Get(id string) (models.Offer, error) {
	obj := Offer{}
	if err := d.db.Take(&obj, "id = ?", id).Error; err != nil {
		return models.Offer{}, fmt.Errorf("take: %w", convertError(err))
	}

	return obj.into(), nil
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

func (d *offerDB) Create(offer *models.Offer) error {
	if offer == nil {
		return ports.ErrOfferIsNil
	}

	o := Offer{}
	o.copy(offer)
	if err := d.db.Create(&o).Error; err != nil {
		return fmt.Errorf("create: %w", convertError(err))
	}

	*offer = o.into()
	return nil
}

func (d *offerDB) Update(offer *models.Offer) error {
	if offer == nil {
		return ports.ErrOfferIsNil
	}

	p := Offer{}
	p.copy(offer)

	if err := d.db.Save(&p).Error; err != nil {
		return fmt.Errorf("save: %w", convertError(err))
	}

	return nil
}

func (d *offerDB) Pagination(first *int, after *string, filter *models.OffersFilter) (models.OffersConnection, error) {
	query := d.filter(d.db.Model(&Offer{}), filter)
	query, err := paginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return models.OffersConnection{}, fmt.Errorf("pagination: %w", err)
	}

	var offers []Offer
	if err := query.Find(&offers).Error; err != nil {
		return models.OffersConnection{}, fmt.Errorf("find: %w", convertError(err))
	}

	if len(offers) == 0 {
		return models.OffersConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.OffersConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(offers) > *first
		offers = offers[:len(offers)-1]
	}

	edges := make([]*models.OffersConnectionEdge, 0, len(offers))

	for _, obj := range offers {
		node := obj.into()
		edges = append(edges, &models.OffersConnectionEdge{
			Cursor: node.ID,
			Node:   &node,
		})
	}

	return models.OffersConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &offers[0].ID,
			EndCursor:   &offers[len(offers)-1].ID,
		},
		Edges: edges,
	}, nil
}
