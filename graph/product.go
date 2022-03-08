package graph

import (
	"auction-back/db"
	"auction-back/graph/model"
	"fmt"

	"gorm.io/gorm"
)

func ProductPagination(query *gorm.DB, first *int, after *string) (*model.ProductsConnection, error) {
	if first != nil {
		if *first < 1 {
			return nil, fmt.Errorf("first must be positive")
		}
		query = query.Limit(*first + 1)
	}

	if after != nil {
		query.Where("id > ?", after)
	}

	var products []db.Product

	result := query.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(products) == 0 {
		return &model.ProductsConnection{
			PageInfo: &model.PageInfo{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     nil,
				EndCursor:       nil,
			},
			Edges: make([]*model.ProductsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(products) > *first
		products = products[:len(products)-1]
	}

	edges := make([]*model.ProductsConnectionEdge, 0, len(products))

	for _, product := range products {
		node, err := (&model.Product{}).From(&product)

		if err != nil {
			return nil, err
		}

		edges = append(edges, &model.ProductsConnectionEdge{
			Cursor: product.ID,
			Node:   node,
		})
	}

	return &model.ProductsConnection{
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &products[0].ID,
			EndCursor:       &products[len(products)-1].ID,
		},
		Edges: edges,
	}, nil
}
