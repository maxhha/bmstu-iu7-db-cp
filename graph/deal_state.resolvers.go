package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *dealStateResolver) Creator(ctx context.Context, obj *models.DealState) (*models.User, error) {
	if obj.CreatorID == nil {
		return nil, nil
	}
	user, err := r.DB.User().Get(*obj.CreatorID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.User().Get: %w", err)
	}

	return &user, nil
}

func (r *dealStateResolver) Offer(ctx context.Context, obj *models.DealState) (*models.Offer, error) {
	offer, err := r.DB.Offer().Get(obj.OfferID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Offer().Get: %w", err)
	}
	return &offer, nil
}

// DealState returns generated.DealStateResolver implementation.
func (r *Resolver) DealState() generated.DealStateResolver { return &dealStateResolver{r} }

type dealStateResolver struct{ *Resolver }
