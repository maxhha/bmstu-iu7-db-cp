package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"context"
	"fmt"
)

func (r *bankAccountResolver) Bank(ctx context.Context, obj *model.BankAccount) (*db.Bank, error) {
	if obj == nil {
		return nil, fmt.Errorf("bank account is nil")
	}

	if err := obj.EnsureFillBank(r.DB); err != nil {
		return nil, err
	}

	return &obj.Bank, nil
}

func (r *bankAccountResolver) Transactions(ctx context.Context, obj *model.BankAccount, first *int, after *string) (*model.TransactionsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userAccountResolver) Bank(ctx context.Context, obj *model.UserAccount) (*db.Bank, error) {
	if obj == nil {
		return nil, fmt.Errorf("user account is nil")
	}

	if err := obj.EnsureFillBank(r.DB); err != nil {
		return nil, err
	}

	return &obj.Bank, nil
}

func (r *userAccountResolver) Transactions(ctx context.Context, obj *model.UserAccount, first *int, after *string) (*model.TransactionsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// BankAccount returns generated.BankAccountResolver implementation.
func (r *Resolver) BankAccount() generated.BankAccountResolver { return &bankAccountResolver{r} }

// UserAccount returns generated.UserAccountResolver implementation.
func (r *Resolver) UserAccount() generated.UserAccountResolver { return &userAccountResolver{r} }

type bankAccountResolver struct{ *Resolver }
type userAccountResolver struct{ *Resolver }
