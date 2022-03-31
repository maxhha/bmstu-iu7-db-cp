package graph

import (
	"auction-back/db"
	"auction-back/graph/model"
	"fmt"

	"gorm.io/gorm"
)

// Creates pagination for user forms
func UserFormPagination(query *gorm.DB, first *int, after *string) (*model.UserFormsConnection, error) {
	query, err := PaginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var objs []db.UserForm
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	if len(objs) == 0 {
		return &model.UserFormsConnection{
			PageInfo: &model.PageInfo{},
			Edges:    make([]*model.UserFormsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*model.UserFormsConnectionEdge, 0, len(objs))

	for _, obj := range objs {
		node := obj

		edges = append(edges, &model.UserFormsConnectionEdge{
			Cursor: obj.ID,
			Node:   &node,
		})
	}

	return &model.UserFormsConnection{
		PageInfo: &model.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, nil
}
