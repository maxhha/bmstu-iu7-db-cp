package graph

import (
	"auction-back/db"
	"auction-back/graph/model"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

func FindAccountPagination(query *gorm.DB, first *int, after *string) ([]db.Account, error) {
	query, err := PaginationByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var objs []db.Account
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	return objs, nil
}

// Creates pagination for accounts
func AccountPagination(query *gorm.DB, first *int, after *string) (*model.AccountsConnection, error) {
	objs, err := FindAccountPagination(query, first, after)

	if err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		return &model.AccountsConnection{
			PageInfo: &model.PageInfo{},
			Edges:    make([]*model.AccountsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*model.AccountsConnectionEdge, 0, len(objs))
	var errors error

	for _, obj := range objs {
		node, err := model.AccountFromDBAccount(obj)

		if err == nil {
			edges = append(edges, &model.AccountsConnectionEdge{
				Cursor: obj.ID,
				Node:   node,
			})
		} else {
			errors = multierror.Append(errors, err)
		}
	}

	return &model.AccountsConnection{
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &objs[0].ID,
			EndCursor:       &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, errors
}

// Creates pagination for accounts
func UserAccountPagination(query *gorm.DB, first *int, after *string) (*model.UserAccountsConnection, error) {
	objs, err := FindAccountPagination(query, first, after)

	if err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		return &model.UserAccountsConnection{
			PageInfo: &model.PageInfo{},
			Edges:    make([]*model.UserAccountsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*model.UserAccountsConnectionEdge, 0, len(objs))
	var errors error

	for _, obj := range objs {

		if obj.Type == db.AccountTypeUser {
			edges = append(edges, &model.UserAccountsConnectionEdge{
				Cursor: obj.ID,
				Node:   &model.UserAccount{Account: obj},
			})
		} else {
			errors = multierror.Append(
				errors,
				fmt.Errorf("unexpected user account type: %s", obj.Type))
		}
	}

	return &model.UserAccountsConnection{
		PageInfo: &model.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, errors
}
