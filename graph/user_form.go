package graph

import (
	"auction-back/db"
	"auction-back/graph/model"
	"fmt"

	"gorm.io/gorm"
)

// Creates pagination for user forms
// query must be ordered: query.Order("created_at desc")
func UserFormPagination(query *gorm.DB, first *int, after *string) (*model.UserFormsConnection, error) {
	if first != nil {
		if *first < 1 {
			return nil, fmt.Errorf("first must be positive")
		}
		query = query.Limit(*first + 1)
	}

	if after != nil {
		query.Where("created_at < ANY(SELECT created_at FROM user_forms u WHERE u.id = ?)", after)
	}

	var objs []db.UserForm

	result := query.Find(&objs)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(objs) == 0 {
		return &model.UserFormsConnection{
			PageInfo: &model.PageInfo{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     nil,
				EndCursor:       nil,
			},
			Edges: make([]*model.UserFormsConnectionEdge, 0),
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
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &objs[0].ID,
			EndCursor:       &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, nil
}
