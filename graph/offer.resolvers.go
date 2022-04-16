package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/models"
	"auction-back/ports"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateOffer(ctx context.Context, input models.CreateOfferInput) (*models.OfferResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	auction, err := r.DB.Auction().Get(input.AuctionID)
	if err != nil {
		return nil, fmt.Errorf("db auction get: %w", err)
	}

	account, err := r.DB.Account().Get(input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("db account get: %w", err)
	}

	if err := isAccountOwner(viewer, account); err != nil {
		return nil, fmt.Errorf("denied: %w", err)
	}

	if !auction.IsStarted() {
		return nil, ErrAuctionIsNotStarted
	}

	if auction.SellerAccountID == nil {
		return nil, fmt.Errorf("seller account id %w", ports.ErrIsNil)
	}

	offer := models.Offer{
		AuctionID: auction.ID,
		UserID:    viewer.ID,
	}

	err = r.Tx(func(tx ports.DB) error {
		if err := tx.Auction().LockShare(&auction); err != nil {
			return fmt.Errorf("db auction lock share: %w", err)
		}

		if err := tx.Account().LockFull(&account); err != nil {
			return fmt.Errorf("db auction lock full: %w", err)
		}

		if auction.IsFinished() {
			return ErrAuctionIsFinished
		}

		// TODO: create fee transaction
		buyTransaction := models.Transaction{
			Type:          models.TransactionTypeBuy,
			Currency:      auction.Currency,
			Amount:        input.Amount, // TODO: calculate by bank transfer table
			AccountFromID: account.ID,
			AccountToID:   *auction.SellerAccountID,
			OfferID:       offer.ID,
		}

		moneys, err := tx.Account().GetAvailableMoney(account)
		if err != nil {
			return fmt.Errorf("db account get available money: %w", err)
		}

		money, exists := moneys[auction.Currency]
		if !exists {
			return fmt.Errorf("%w: %s", ErrNoCurrency, auction.Currency)
		}

		if money.Amount.LessThan(buyTransaction.Amount) {
			return fmt.Errorf(
				"%w: need %s but have %s",
				ErrNotAnoughMoney,
				buyTransaction.Amount.String(),
				money.Amount.String(),
			)
		}

		if err := tx.Offer().Create(&offer); err != nil {
			return fmt.Errorf("db offer create: %w", err)
		}

		buyTransaction.OfferID = offer.ID

		if err := tx.Transaction().Create(&buyTransaction); err != nil {
			return fmt.Errorf("db transaction create: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("tx: %w", err)
	}

	return &models.OfferResult{Offer: &offer}, nil
}

func (r *offerResolver) User(ctx context.Context, obj *models.Offer) (*models.User, error) {
	if obj == nil {
		return nil, nil
	}

	user, err := r.DB.User().Get(obj.UserID)
	if err != nil {
		return nil, fmt.Errorf("db user get: %w", err)
	}

	return &user, nil
}

func (r *offerResolver) Auction(ctx context.Context, obj *models.Offer) (*models.Auction, error) {
	if obj == nil {
		return nil, nil
	}

	auction, err := r.DB.Auction().Get(obj.AuctionID)
	if err != nil {
		return nil, fmt.Errorf("db user get: %w", err)
	}

	return &auction, nil
}

func (r *offerResolver) Moneys(ctx context.Context, obj *models.Offer) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) Transactions(ctx context.Context, obj *models.Offer) ([]*models.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Offers(ctx context.Context, first *int, after *string, filter *models.OffersFilter) (*models.OffersConnection, error) {
	connection, err := r.DB.Offer().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db offer pagination: %w", err)
	}

	return &connection, nil
}

// Offer returns generated.OfferResolver implementation.
func (r *Resolver) Offer() generated.OfferResolver { return &offerResolver{r} }

type offerResolver struct{ *Resolver }
