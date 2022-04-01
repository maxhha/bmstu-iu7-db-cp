package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"context"
)

func (r *bankAccountResolver) Bank(ctx context.Context, obj *model.BankAccount) (*db.Bank, error) {
	return r.Account().Bank(ctx, obj.AccountPtr())
}

func (r *bankAccountResolver) Transactions(ctx context.Context, obj *model.BankAccount, first *int, after *string) (*model.TransactionsConnection, error) {
	return r.Account().Transactions(ctx, obj.AccountPtr(), first, after)
}

func (r *userAccountResolver) Bank(ctx context.Context, obj *model.UserAccount) (*db.Bank, error) {
	return r.Account().Bank(ctx, obj.AccountPtr())
}

func (r *userAccountResolver) Transactions(ctx context.Context, obj *model.UserAccount, first *int, after *string) (*model.TransactionsConnection, error) {
	return r.Account().Transactions(ctx, obj.AccountPtr(), first, after)
}

// BankAccount returns generated.BankAccountResolver implementation.
func (r *Resolver) BankAccount() generated.BankAccountResolver { return &bankAccountResolver{r} }

// UserAccount returns generated.UserAccountResolver implementation.
func (r *Resolver) UserAccount() generated.UserAccountResolver { return &userAccountResolver{r} }

type bankAccountResolver struct{ *Resolver }
type userAccountResolver struct{ *Resolver }
