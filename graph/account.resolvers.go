package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *bankAccountResolver) Bank(ctx context.Context, obj *models.BankAccount) (*models.Bank, error) {
	return r.Account().Bank(ctx, obj.AccountPtr())
}

func (r *bankAccountResolver) Available(ctx context.Context, obj *models.BankAccount) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *bankAccountResolver) Blocked(ctx context.Context, obj *models.BankAccount) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *bankAccountResolver) Transactions(ctx context.Context, obj *models.BankAccount, first *int, after *string, filter *models.TransactionsFilter) (*models.TransactionsConnection, error) {
	return r.Account().Transactions(ctx, obj.AccountPtr(), first, after)
}

func (r *queryResolver) Accounts(ctx context.Context, first *int, after *string, filter *models.AccountsFilter) (*models.AccountsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userAccountResolver) Bank(ctx context.Context, obj *models.UserAccount) (*models.Bank, error) {
	return r.Account().Bank(ctx, obj.AccountPtr())
}

func (r *userAccountResolver) Available(ctx context.Context, obj *models.UserAccount) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userAccountResolver) Blocked(ctx context.Context, obj *models.UserAccount) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userAccountResolver) Transactions(ctx context.Context, obj *models.UserAccount, first *int, after *string, filter *models.TransactionsFilter) (*models.TransactionsConnection, error) {
	return r.Account().Transactions(ctx, obj.AccountPtr(), first, after)
}

func (r *userAccountResolver) User(ctx context.Context, obj *models.UserAccount) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// BankAccount returns generated.BankAccountResolver implementation.
func (r *Resolver) BankAccount() generated.BankAccountResolver { return &bankAccountResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// UserAccount returns generated.UserAccountResolver implementation.
func (r *Resolver) UserAccount() generated.UserAccountResolver { return &userAccountResolver{r} }

type bankAccountResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userAccountResolver struct{ *Resolver }
