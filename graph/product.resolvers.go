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
		Owner:       *viewer,
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
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	product := db.Product{}

	result := db.DB.Take(&product, "id = ?", input.ProductID)

	if result.Error != nil {
		return nil, fmt.Errorf("db take: %w", result.Error)
	}

	if product.OwnerID != viewer.ID {
		return nil, fmt.Errorf("viewer is not owner")
	}

	product.IsOnMarket = true

	result = db.DB.Save(&product)

	if result.Error != nil {
		return nil, fmt.Errorf("db save: %w", result.Error)
	}

	p, err := (&model.Product{}).From(&product)

	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return &model.OfferProductResult{Product: p}, nil
}

func (r *mutationResolver) TakeOffProduct(ctx context.Context, input model.TakeOffProductInput) (*model.TakeOffProductResult, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	product := db.Product{}

	result := db.DB.Take(&product, "id = ?", input.ProductID)

	if result.Error != nil {
		return nil, fmt.Errorf("db take: %w", result.Error)
	}

	if product.OwnerID != viewer.ID {
		return nil, fmt.Errorf("viewer is not owner")
	}

	product.IsOnMarket = true

	result = db.DB.Save(&product)

	if result.Error != nil {
		return nil, fmt.Errorf("db save: %w", result.Error)
	}

	p, err := (&model.Product{}).From(&product)

	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return &model.TakeOffProductResult{Product: p}, nil
}

func (r *mutationResolver) SellProduct(ctx context.Context, input model.SellProductInput) (*model.SellProductResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Owner(ctx context.Context, obj *model.Product) (*model.User, error) {
	if obj.DB.Owner.ID == obj.DB.OwnerID {
		return (&model.User{}).From(&obj.DB.Owner)
	}

	owner := db.User{}
	result := db.DB.Take(&owner, "id = ?", obj.DB.OwnerID)

	if result.Error != nil {
		return nil, fmt.Errorf("db take: %w", result.Error)
	}

	return (&model.User{}).From(&owner)
}

func (r *productResolver) Offers(ctx context.Context, obj *model.Product, first *int, after *string) (*model.OffersConnection, error) {
	query := db.DB.Where("product_id = ?", obj.ID).Order("id")

	return OfferPagination(query, first, after)
}

func (r *queryResolver) MarketProducts(ctx context.Context, first *int, after *string) (*model.ProductConnection, error) {
	query := db.DB.Where("is_on_market = true").Order("id")

	return ProductPagination(query, first, after)
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
