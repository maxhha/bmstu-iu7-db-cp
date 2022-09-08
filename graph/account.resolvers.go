package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *accountResolver) User(ctx context.Context, obj *models.Account) (*models.User, error) {
	user, err := r.DB.User().Get(obj.UserID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.User(): %w", err)
	}

	return &user, nil
}

func (r *accountResolver) NominalAccount(ctx context.Context, obj *models.Account) (*models.NominalAccount, error) {
	acc, err := r.DB.NominalAccount().Get(obj.NominalAccountID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.NominalAccount().Get: %w", err)
	}

	return &acc, nil
}

func (r *accountResolver) Available(ctx context.Context, obj *models.Account) ([]*models.Money, error) {
	money, err := r.DB.Account().GetAvailableMoney(*obj)
	if err != nil {
		return nil, fmt.Errorf("DB.Account().GetAvailableMoney: %w", err)
	}

	return moneyMapToArray(money), nil
}

func (r *accountResolver) Blocked(ctx context.Context, obj *models.Account) ([]*models.Money, error) {
	money, err := r.DB.Account().GetBlockedMoney(*obj)
	if err != nil {
		return nil, fmt.Errorf("DB.Account().GetBlockedMoney: %w", err)
	}

	return moneyMapToArray(money), nil
}

func (r *accountResolver) Transactions(ctx context.Context, obj *models.Account, first *int, after *string, filter *models.TransactionsFilter) (*models.TransactionsConnection, error) {
	if filter == nil {
		filter = &models.TransactionsFilter{}
	}

	filter.AccountIDs = []string{obj.ID}

	trs, err := r.DB.Transaction().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Transaction().Pagination: %w", err)
	}

	return &trs, nil
}

func (r *mutationResolver) CreateAccount(ctx context.Context, input models.CreateAccountInput) (*models.AccountResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	account, err := r.BankPort.CreateAccount(viewer.ID, input.NominalAccountID)
	return &models.AccountResult{
		Account: &account,
	}, nil
}

func (r *queryResolver) Accounts(ctx context.Context, first *int, after *string, filter *models.AccountsFilter) (*models.AccountsConnection, error) {
	conn, err := r.DB.Account().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Account().Pagination: %w", err)
	}
	return &conn, nil
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type accountResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
