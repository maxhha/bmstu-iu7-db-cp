package database

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"
	"time"

	"gorm.io/gorm"
)

//go:generate go run ../../codegen/gormdbops/main.go --out deal_state_gen.go --model DealState --methods Get,Create

type DealState struct {
	ID        string               `gorm:"default:generated();"`
	State     models.DealStateEnum `gorm:"default:'TRANSFERRING_MONEY';"`
	CreatorID *string
	OfferID   string
	Comment   *string
	CreatedAt time.Time
}

func (a *DealState) into() models.DealState {
	return models.DealState{
		ID:        a.ID,
		State:     a.State,
		CreatorID: a.CreatorID,
		OfferID:   a.OfferID,
		Comment:   a.Comment,
		CreatedAt: a.CreatedAt,
	}
}

func (a *DealState) copy(obj *models.DealState) {
	if obj == nil {
		return
	}

	a.ID = obj.ID
	a.State = obj.State
	a.CreatorID = obj.CreatorID
	a.OfferID = obj.OfferID
	a.Comment = obj.Comment
	a.CreatedAt = obj.CreatedAt
}

func (d *dealStateDB) filter(query *gorm.DB, config *models.DealStateFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.CreatorIDs) > 0 {
		query = query.Where("creator_id IN ?", config.CreatorIDs)
	}

	if len(config.OfferIDs) > 0 {
		query = query.Where("offer_id IN ?", config.OfferIDs)
	}

	return query
}

func (d *dealStateDB) Find(config ports.DealStateFindConfig) ([]models.DealState, error) {
	query := d.filter(d.db, config.Filter)
	query = query.Order("created_at DESC")

	if config.Limit > 0 {
		query = query.Limit(config.Limit)
	}

	var objs []DealState
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", convertError(err))
	}

	arr := make([]models.DealState, 0, len(objs))
	for _, obj := range objs {
		arr = append(arr, obj.into())
	}

	return arr, nil
}

func (d *dealStateDB) GetLast(offerId string) (models.DealState, error) {
	query := d.filter(d.db, &models.DealStateFilter{
		OfferIDs: []string{offerId},
	})
	query = query.Order("created_at DESC")

	var obj DealState
	if err := query.Take(&obj).Error; err != nil {
		return models.DealState{}, fmt.Errorf("take: %w", convertError(err))
	}

	return obj.into(), nil
}
