package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
)

func (r *moneyResolver) Amount(ctx context.Context, obj *models.Money) (float64, error) {
	if obj == nil {
		return 0.0, nil
	}

	f, exact := obj.Amount.Float64()
	if !exact {
		return f, ErrFloatIsNotExact
	}

	return f, nil
}

// Money returns generated.MoneyResolver implementation.
func (r *Resolver) Money() generated.MoneyResolver { return &moneyResolver{r} }

type moneyResolver struct{ *Resolver }
