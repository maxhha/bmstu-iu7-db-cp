package graph

import (
	"auction-back/models"
	"fmt"

	"gorm.io/gorm"
)

func OfferPagination(query *gorm.DB, first *int, after *string) (*models.OffersConnection, error) {
	if first != nil {
		if *first < 1 {
			return nil, fmt.Errorf("first must be positive")
		}
		query = query.Limit(*first + 1)
	}

	if after != nil {
		query.Where("id > ?", after)
	}

	var offers []models.Offer

	result := query.Find(&offers)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(offers) == 0 {
		return &models.OffersConnection{
			PageInfo: &models.PageInfo{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     nil,
				EndCursor:       nil,
			},
			Edges: make([]*models.OffersConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(offers) > *first
		offers = offers[:len(offers)-1]
	}

	edges := make([]*models.OffersConnectionEdge, 0, len(offers))

	for _, offer := range offers {
		edges = append(edges, &models.OffersConnectionEdge{
			Cursor: offer.ID,
			Node:   &offer,
		})
	}

	return &models.OffersConnection{
		PageInfo: &models.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &offers[0].ID,
			EndCursor:       &offers[len(offers)-1].ID,
		},
		Edges: edges,
	}, nil
}
