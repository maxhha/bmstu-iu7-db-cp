package graph

import (
	"auction-back/models"
	"fmt"

	"gorm.io/gorm"
)

func ProductPagination(query *gorm.DB, first *int, after *string) (*models.ProductsConnection, error) {
	query, err := PaginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var products []models.Product
	result := query.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(products) == 0 {
		return &models.ProductsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.ProductsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(products) > *first
		products = products[:len(products)-1]
	}

	edges := make([]*models.ProductsConnectionEdge, 0, len(products))

	for _, node := range products {
		edges = append(edges, &models.ProductsConnectionEdge{
			Cursor: node.ID,
			Node:   &node,
		})
	}

	return &models.ProductsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &products[0].ID,
			EndCursor:   &products[len(products)-1].ID,
		},
		Edges: edges,
	}, nil
}
