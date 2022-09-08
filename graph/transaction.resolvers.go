package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *queryResolver) Transactions(ctx context.Context, first *int, after *string, filter *models.TransactionsFilter) (*models.TransactionsConnection, error) {
	conn, err := r.DB.Transaction().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Transaction().Pagination: %w", err)
	}
	return &conn, nil
}

func (r *transactionResolver) AccountFrom(ctx context.Context, obj *models.Transaction) (*models.Account, error) {
	if obj.AccountFromID == nil {
		return nil, nil
	}
	acc, err := r.DB.Account().Get(*obj.AccountFromID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Account().Get: %w", err)
	}

	return &acc, nil
}

func (r *transactionResolver) AccountTo(ctx context.Context, obj *models.Transaction) (*models.Account, error) {
	if obj.AccountToID == nil {
		return nil, nil
	}
	acc, err := r.DB.Account().Get(*obj.AccountToID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Account().Get: %w", err)
	}

	return &acc, nil
}

func (r *transactionResolver) Offer(ctx context.Context, obj *models.Transaction) (*models.Offer, error) {
	if obj.OfferID == nil {
		return nil, nil
	}
	off, err := r.DB.Offer().Get(*obj.OfferID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Offer().Get: %w", err)
	}
	return &off, nil
}

// Transaction returns generated.TransactionResolver implementation.
func (r *Resolver) Transaction() generated.TransactionResolver { return &transactionResolver{r} }

type transactionResolver struct{ *Resolver }
