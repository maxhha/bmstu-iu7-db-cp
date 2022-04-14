package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
	"time"
)

func (r *auctionResolver) Product(ctx context.Context, obj *models.Auction) (*models.Product, error) {
	if obj == nil {
		return nil, nil
	}

	product, err := r.DB.Product().Get(obj.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product db get: %w", err)
	}

	return &product, nil
}

func (r *auctionResolver) Seller(ctx context.Context, obj *models.Auction) (*models.User, error) {
	if obj == nil {
		return nil, nil
	}

	user, err := r.DB.User().Get(obj.SellerID)
	if err != nil {
		return nil, fmt.Errorf("user db get: %w", err)
	}

	return &user, nil
}

func (r *auctionResolver) Buyer(ctx context.Context, obj *models.Auction) (*models.User, error) {
	if obj == nil || obj.BuyerID == nil {
		return nil, nil
	}

	user, err := r.DB.User().Get(*obj.BuyerID)
	if err != nil {
		return nil, fmt.Errorf("user db get: %w", err)
	}

	return &user, nil
}

func (r *mutationResolver) CreateAuction(ctx context.Context, input models.ProductInput) (*models.AuctionResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	product, err := r.DB.Product().Get(input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	if err := isProductOwner(r.DB, viewer, product); err != nil {
		return nil, err
	}

	if product.State != models.ProductStateApproved {
		return nil, fmt.Errorf("product state is not %s", models.ProductStateApproved)
	}

	auction := models.Auction{
		ProductID: product.ID,
		SellerID:  viewer.ID,
	}
	if err := r.DB.Auction().Create(&auction); err != nil {
		return nil, fmt.Errorf("db create: %w", err)
	}

	return &models.AuctionResult{
		Auction: &auction,
	}, nil
}

func (r *mutationResolver) UpdateAuction(ctx context.Context, input models.UpdateAuctionInput) (*models.AuctionResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	auction, err := r.DB.Auction().Get(input.AuctionID)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	if err := isAuctionOwner(viewer, auction); err != nil {
		return nil, err
	}

	if !auction.IsEditable() {
		return nil, ErrNotEditable
	}

	auction.MinMoney = input.MinMoney.IntoPtr()
	auction.ScheduledStartAt = input.ScheduledStartAt
	auction.ScheduledFinishAt = input.ScheduledFinishAt

	if err := r.DB.Auction().Update(&auction); err != nil {
		return nil, fmt.Errorf("db auction update: %w", err)
	}

	return &models.AuctionResult{Auction: &auction}, nil
}

func (r *mutationResolver) StartAuction(ctx context.Context, input models.AuctionInput) (*models.AuctionResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	auction, err := r.DB.Auction().Get(input.AuctionID)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	if err := isAuctionOwner(viewer, auction); err != nil {
		return nil, err
	}

	if !auction.IsEditable() {
		return nil, ErrNotEditable
	}

	now := time.Now()
	auction.ScheduledStartAt = &now
	auction.State = models.AuctionStateStarted

	if err := r.DB.Auction().Update(&auction); err != nil {
		return nil, fmt.Errorf("db auction update: %w", err)
	}

	return &models.AuctionResult{Auction: &auction}, nil
}

func (r *queryResolver) Auctions(ctx context.Context, first *int, after *string, filter *models.AuctionsFilter) (*models.AuctionsConnection, error) {
	connection, err := r.DB.Auction().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}

	return &connection, nil
}

// Auction returns generated.AuctionResolver implementation.
func (r *Resolver) Auction() generated.AuctionResolver { return &auctionResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type auctionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
