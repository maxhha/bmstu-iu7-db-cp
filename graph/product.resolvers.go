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

	"github.com/teris-io/shortid"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.CreateProductResult, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	id, err := shortid.Generate()

	if err != nil {
		return nil, fmt.Errorf("shortid: %w", err)
	}

	product := db.Product{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     viewer.ID,
	}

	result := db.DB.Create(&product)

	if result.Error != nil {
		return nil, fmt.Errorf("db create: %w", result.Error)
	}

	p, err := (&model.Product{}).From(&product)

	if err != nil {
		return nil, fmt.Errorf("convert: %w", result.Error)
	}

	return &model.CreateProductResult{
		Product: p,
	}, nil
}

func (r *mutationResolver) OfferProduct(ctx context.Context, input model.OfferProductInput) (*model.OfferProductResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TakeOffProduct(ctx context.Context, input model.TakeOffProductInput) (*model.TakeOffProductResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SellProduct(ctx context.Context, input model.SellProductInput) (*model.SellProductResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Owner(ctx context.Context, obj *model.Product) (*model.User, error) {
	owner := db.User{}
	result := db.DB.Take(&owner, "id = ?", obj.DB.OwnerID)

	if result.Error != nil {
		return nil, fmt.Errorf("db take: %w", result.Error)
	}

	return (&model.User{}).From(&owner)
}

func (r *productResolver) Offers(ctx context.Context, obj *model.Product, first *int, after *string) (*model.OffersConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *productResolver) IsOnMarket(ctx context.Context, obj *model.Product) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}
