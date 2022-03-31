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

func (r *bankResolver) Account(ctx context.Context, obj *db.Bank) (*model.BankAccount, error) {
	if obj == nil {
		return nil, fmt.Errorf("bank is nil")
	}

	account := db.Account{}

	if err := r.DB.Take(&account, "type = ? AND bank_id = ?", db.AccountTypeBank, obj.ID).Error; err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	return &model.BankAccount{Account: account}, nil
}

// Bank returns generated.BankResolver implementation.
func (r *Resolver) Bank() generated.BankResolver { return &bankResolver{r} }

type bankResolver struct{ *Resolver }
