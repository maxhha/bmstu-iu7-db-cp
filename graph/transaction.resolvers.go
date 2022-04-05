package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *transactionResolver) ID(ctx context.Context, obj *models.Transaction) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) State(ctx context.Context, obj *models.Transaction) (models.TransactionStateEnum, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) Type(ctx context.Context, obj *models.Transaction) (models.TransactionTypeEnum, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) Currency(ctx context.Context, obj *models.Transaction) (models.CurrencyEnum, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) Amount(ctx context.Context, obj *models.Transaction) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) AccountFrom(ctx context.Context, obj *models.Transaction) (models.AccountInterface, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) AccountTo(ctx context.Context, obj *models.Transaction) (models.AccountInterface, error) {
	panic(fmt.Errorf("not implemented"))
}

// Transaction returns generated.TransactionResolver implementation.
func (r *Resolver) Transaction() generated.TransactionResolver { return &transactionResolver{r} }

type transactionResolver struct{ *Resolver }
