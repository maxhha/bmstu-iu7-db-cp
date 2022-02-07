package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateOffer(ctx context.Context, input model.CreateOfferInput) (*model.CreateOfferResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveOffer(ctx context.Context, input model.RemoveOfferInput) (*model.RemoveOfferResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
