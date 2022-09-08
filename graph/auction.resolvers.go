package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/models"
	"auction-back/ports"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
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

func (r *auctionResolver) SellerAccount(ctx context.Context, obj *models.Auction) (*models.Account, error) {
	if obj.SellerAccountID == nil {
		return nil, nil
	}
	acc, err := r.DB.Account().Get(*obj.SellerAccountID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Account().Get: %w", err)
	}

	return &acc, nil
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

func (r *auctionResolver) TopOffer(ctx context.Context, obj *models.Auction) (*models.Offer, error) {
	offer, err := r.DB.Auction().GetTopOffer(*obj)
	if err != nil {
		if errors.Is(err, ports.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("r.DB.Auction().GetTopOffer: %w", err)
	}

	return &offer, nil
}

func (r *auctionResolver) Offers(ctx context.Context, obj *models.Auction, first *int, after *string, filter *models.OffersFilter) (*models.OffersConnection, error) {
	if filter == nil {
		filter = &models.OffersFilter{}
	}

	if len(filter.AuctionIDs) > 0 {
		return nil, fmt.Errorf("auctionIDs must be empty")
	}

	filter.AuctionIDs = []string{obj.ID}
	connection, err := r.DB.Offer().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Offer().Pagination: %w", err)
	}

	return &connection, nil
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

	if err := IsProductOwner(r.DB, viewer, product); err != nil {
		return nil, err
	}

	if product.State != models.ProductStateApproved {
		return nil, fmt.Errorf("product state is not %s", models.ProductStateApproved)
	}

	auction, err := r.DB.Auction().Take(ports.AuctionTakeConfig{
		Filter: &models.AuctionsFilter{
			SellerIDs:  []string{viewer.ID},
			ProductIDs: []string{product.ID},
		},
	})
	if err == nil {
		if auction.State != models.AuctionStateFailed {
			return nil, ErrAlreadyExists
		}
	} else if !errors.Is(err, ports.ErrRecordNotFound) {
		return nil, fmt.Errorf("db auction take: %w", err)
	}

	form, err := r.DB.User().LastApprovedUserForm(viewer)
	if err != nil {
		return nil, fmt.Errorf("db user last approved form: %w", err)
	}
	if form.Currency == nil {
		return nil, ErrCurrencyIsNil
	}

	auction = models.Auction{
		ProductID: product.ID,
		SellerID:  viewer.ID,
		Currency:  *form.Currency,
	}
	if err := r.DB.Auction().Create(&auction); err != nil {
		return nil, fmt.Errorf("db auction create: %w", err)
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

	if err := IsAuctionOwner(viewer, auction); err != nil {
		return nil, err
	}

	if !auction.IsEditable() {
		return nil, ErrNotEditable
	}

	auction.Currency = input.Currency
	auction.SellerAccountID = input.SellerAccountID
	auction.MinAmount = input.MinAmount
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

	if err := IsAuctionOwner(viewer, auction); err != nil {
		return nil, err
	}

	if !auction.IsEditable() {
		return nil, ErrNotEditable
	}

	now := time.Now().UTC()
	auction.StartedAt = &now
	auction.State = models.AuctionStateStarted

	if err := r.DB.Auction().Update(&auction); err != nil {
		return nil, fmt.Errorf("db auction update: %w", err)
	}

	return &models.AuctionResult{Auction: &auction}, nil
}

func (r *mutationResolver) ApproveOwnerAcceptedMoney(ctx context.Context, input models.AuctionInput) (*models.AuctionResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	auction, err := r.DB.Auction().Get(input.AuctionID)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	if err := IsAuctionOwner(viewer, auction); err != nil {
		return nil, err
	}

	offer, err := r.DB.Offer().Take(ports.OfferTakeConfig{
		Filter: &models.OffersFilter{
			States:     []models.OfferState{models.OfferStateAccepted},
			AuctionIDs: []string{input.AuctionID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("r.DB.Offer().Take: %w", err)
	}

	err = r.DealerPort.OwnerAccept(offer.ID, &viewer.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("r.DealerPort.OwnerAccept: %w", err)
	}

	return &models.AuctionResult{
		Auction: &auction,
	}, nil
}

func (r *mutationResolver) ApproveReceiveProduct(ctx context.Context, input models.AuctionInput) (bool, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return false, err
	}

	offer, err := r.DB.Offer().Take(ports.OfferTakeConfig{Filter: &models.OffersFilter{
		AuctionIDs: []string{input.AuctionID},
		UserIDs:    []string{viewer.ID},
		States:     []models.OfferState{models.OfferStateAccepted},
	}})
	if err != nil {
		return false, fmt.Errorf("r.DB.Offer().Find: %w", err)
	}

	err = r.DealerPort.BuyerAccept(offer.ID, &viewer.ID)
	if err != nil {
		return false, fmt.Errorf("r.DealerPort.BuyerAccept: %w", err)
	}

	return true, nil
}

func (r *queryResolver) Auctions(ctx context.Context, first *int, after *string, filter *models.AuctionsFilter) (*models.AuctionsConnection, error) {
	connection, err := r.DB.Auction().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}

	return &connection, nil
}

func (r *queryResolver) RunningAuctions(ctx context.Context, first *int, after *string, filter *models.AuctionsFilter) (*models.AuctionsConnection, error) {
	if filter == nil {
		filter = &models.AuctionsFilter{}
	}

	var errs error

	if len(filter.States) > 0 {
		errs = multierror.Append(errs, fmt.Errorf("states must be empty"))
	}

	if filter.StartedAt == nil {
		filter.StartedAt = &models.DateTimeRange{}
	}

	if filter.StartedAt.To != nil {
		errs = multierror.Append(errs, fmt.Errorf("startedAt.to must be empty"))
	}

	if errs != nil {
		return nil, errs
	}

	filter.States = []models.AuctionState{models.AuctionStateStarted}
	now := time.Now()
	filter.StartedAt.To = &now

	connection, err := r.DB.Auction().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Auction().Pagination: %w", err)
	}

	return &connection, nil
}

func (r *queryResolver) ScheduledAuctions(ctx context.Context, first *int, after *string, filter *models.AuctionsFilter) (*models.AuctionsConnection, error) {
	if filter == nil {
		filter = &models.AuctionsFilter{}
	}

	var errs error

	if len(filter.States) > 0 {
		errs = multierror.Append(errs, fmt.Errorf("states must be empty"))
	}

	if filter.ScheduledStartAt != nil {
		errs = multierror.Append(errs, fmt.Errorf("scheduledStartAt must be empty"))
	}

	if errs != nil {
		return nil, errs
	}

	filter.States = []models.AuctionState{models.AuctionStateCreated}
	now := time.Now()
	filter.ScheduledStartAt = &models.DateTimeRange{
		From: &now,
	}

	connection, err := r.DB.Auction().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Auction().Pagination: %w", err)
	}

	return &connection, nil
}

// Auction returns generated.AuctionResolver implementation.
func (r *Resolver) Auction() generated.AuctionResolver { return &auctionResolver{r} }

type auctionResolver struct{ *Resolver }
