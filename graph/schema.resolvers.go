package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"context"
	"fmt"
)

func (r *mutationResolver) Register(ctx context.Context) (*model.RegisterResult, error) {
	user := db.User{}

	result := db.DB.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &model.RegisterResult{
		Token: fmt.Sprintf("%v", user.ID),
	}, nil
}

func (r *mutationResolver) IncreaseBalance(ctx context.Context, input model.IncreaseBalanceInput) (*model.IncreaseBalanceResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.CreateProductResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Viewer(ctx context.Context) (*model.User, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, nil
	}

	return &model.User{
		ID: fmt.Sprintf("%v", viewer.ID),
		Balance: &model.Balance{
			Available: float64(viewer.Available),
			Blocked:   0,
		},
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
