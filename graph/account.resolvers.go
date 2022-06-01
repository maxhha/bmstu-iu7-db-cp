package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *accountResolver) User(ctx context.Context, obj *models.Account) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *accountResolver) NominalAccount(ctx context.Context, obj *models.Account) (*models.NominalAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *accountResolver) Available(ctx context.Context, obj *models.Account) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *accountResolver) Blocked(ctx context.Context, obj *models.Account) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *accountResolver) Transactions(ctx context.Context, obj *models.Account, first *int, after *string, filter *models.TransactionsFilter) (*models.TransactionsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateAccount(ctx context.Context, input models.CreateAccountInput) (*models.AccountResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Accounts(ctx context.Context, first *int, after *string, filter *models.AccountsFilter) (*models.AccountsConnection, error) {
	panic(fmt.Errorf("not implemented"))
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
