package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *bankResolver) Account(ctx context.Context, obj *models.Bank) (*models.BankAccount, error) {
	if obj == nil {
		return nil, fmt.Errorf("bank is nil")
	}

	account, err := r.DB.Bank().GetAccount(*obj)
	if err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	return &account, nil
}

// Bank returns generated.BankResolver implementation.
func (r *Resolver) Bank() generated.BankResolver { return &bankResolver{r} }

type bankResolver struct{ *Resolver }
