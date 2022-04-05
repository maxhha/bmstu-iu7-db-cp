package graph

import (
	"auction-back/models"
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

func FindAccountPagination(query *gorm.DB, first *int, after *string) ([]models.Account, error) {
	query, err := PaginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var objs []models.Account
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	return objs, nil
}

// Creates pagination for accounts
func AccountPagination(query *gorm.DB, first *int, after *string) (*models.AccountsConnection, error) {
	objs, err := FindAccountPagination(query, first, after)

	if err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		return &models.AccountsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.AccountsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.AccountsConnectionEdge, 0, len(objs))
	var errors error

	for _, obj := range objs {
		node, err := obj.ConcreteType()

		if err == nil {
			edges = append(edges, &models.AccountsConnectionEdge{
				Cursor: obj.ID,
				Node:   node,
			})
		} else {
			errors = multierror.Append(errors, err)
		}
	}

	return &models.AccountsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &objs[0].ID,
			EndCursor:       &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, errors
}

// Creates pagination for accounts
func UserAccountPagination(query *gorm.DB, first *int, after *string) (*models.UserAccountsConnection, error) {
	objs, err := FindAccountPagination(query, first, after)

	if err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		return &models.UserAccountsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.UserAccountsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.UserAccountsConnectionEdge, 0, len(objs))
	var errors error

	for _, obj := range objs {

		if obj.Type == models.AccountTypeUser {
			edges = append(edges, &models.UserAccountsConnectionEdge{
				Cursor: obj.ID,
				Node:   &models.UserAccount{Account: obj},
			})
		} else {
			errors = multierror.Append(
				errors,
				fmt.Errorf("unexpected user account type: %s", obj.Type))
		}
	}

	return &models.UserAccountsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, errors
}

func (r *accountResolver) Bank(ctx context.Context, obj *models.Account) (*models.Bank, error) {
	if obj == nil {
		return nil, fmt.Errorf("account is nil")
	}

	if err := obj.EnsureFillBank(r.DB); err != nil {
		return nil, err
	}

	return &obj.Bank, nil
}

func (r *accountResolver) Transactions(ctx context.Context, obj *models.Account, first *int, after *string) (*models.TransactionsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Account() *accountResolver { return &accountResolver{r} }

type accountResolver struct{ *Resolver }
