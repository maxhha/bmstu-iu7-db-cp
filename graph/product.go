package graph

import (
	"auction-back/db"
	"auction-back/graph/model"
	"fmt"

	"gorm.io/gorm"
)

func ProductPagination(query *gorm.DB, first *int, after *string) (*model.ProductsConnection, error) {
	query, err := PaginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var products []db.Product
	result := query.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(products) == 0 {
		return &model.ProductsConnection{
			PageInfo: &model.PageInfo{},
			Edges:    make([]*model.ProductsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(products) > *first
		products = products[:len(products)-1]
	}

	edges := make([]*model.ProductsConnectionEdge, 0, len(products))

	for _, node := range products {
		edges = append(edges, &model.ProductsConnectionEdge{
			Cursor: node.ID,
			Node:   &node,
		})
	}

	return &model.ProductsConnection{
		PageInfo: &model.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &products[0].ID,
			EndCursor:   &products[len(products)-1].ID,
		},
		Edges: edges,
	}, nil
}
