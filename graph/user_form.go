package graph

import (
	"auction-back/models"
	"fmt"

	"gorm.io/gorm"
)

// Creates pagination for user forms
func UserFormPagination(query *gorm.DB, first *int, after *string) (*models.UserFormsConnection, error) {
	query, err := PaginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var objs []models.UserForm
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	if len(objs) == 0 {
		return &models.UserFormsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.UserFormsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.UserFormsConnectionEdge, 0, len(objs))

	for _, obj := range objs {
		node := obj

		edges = append(edges, &models.UserFormsConnectionEdge{
			Cursor: obj.ID,
			Node:   &node,
		})
	}

	return &models.UserFormsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, nil
}
