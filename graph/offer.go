package graph

import (
	"auction-back/db"
	"auction-back/graph/model"
	"fmt"

	"gorm.io/gorm"
)

func OfferPagination(query *gorm.DB, first *int, after *string) (*model.OffersConnection, error) {
	if first != nil {
		if *first < 1 {
			return nil, fmt.Errorf("first must be positive")
		}
		query = query.Limit(*first + 1)
	}

	if after != nil {
		query.Where("id > ?", after)
	}

	var offers []db.Offer

	result := query.Find(&offers)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(offers) == 0 {
		return &model.OffersConnection{
			PageInfo: &model.PageInfo{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     nil,
				EndCursor:       nil,
			},
			Edges: make([]*model.OffersConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(offers) > *first
		offers = offers[:len(offers)-1]
	}

	edges := make([]*model.OffersConnectionEdge, 0, len(offers))

	for _, offer := range offers {
		node, err := (&model.Offer{}).From(&offer)

		if err != nil {
			return nil, err
		}

		edges = append(edges, &model.OffersConnectionEdge{
			Cursor: offer.ID,
			Node:   node,
		})
	}

	return &model.OffersConnection{
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &offers[0].ID,
			EndCursor:       &offers[len(offers)-1].ID,
		},
		Edges: edges,
	}, nil
}
