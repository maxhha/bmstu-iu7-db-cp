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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *mutationResolver) CreateOffer(ctx context.Context, input model.CreateOfferInput) (*model.CreateOfferResult, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	var product db.Product

	result := db.DB.Take(&product, "id = ?", input.ProductID)

	if result.Error != nil {
		return nil, fmt.Errorf("db take product: %w", result.Error)
	}

	if !product.IsOnMarket {
		return nil, fmt.Errorf("product is not on market")
	}

	id, err := shortid.Generate()

	if err != nil {
		return nil, fmt.Errorf("shortid: %w", err)
	}

	offer := db.Offer{
		ID:         id,
		Amount:     input.Amount,
		ProductID:  input.ProductID,
		Product:    product,
		ConsumerID: viewer.ID,
		Consumer:   *viewer,
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(viewer, "id = ?", viewer.ID).Error; err != nil {
			return err
		}

		if viewer.Available < input.Amount {
			return fmt.Errorf("not enough amount")
		}

		// FIXME: presicion
		viewer.Available = viewer.Available - input.Amount

		if err := tx.Save(&viewer).Error; err != nil {
			return err
		}

		if err := tx.Save(&offer).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("db save: %w", err)
	}

	o, err := (&model.Offer{}).From(&offer)

	if err != nil {
		return nil, err
	}

	return &model.CreateOfferResult{
		Offer: o,
	}, nil
}

func (r *mutationResolver) RemoveOffer(ctx context.Context, input model.RemoveOfferInput) (*model.RemoveOfferResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) Consumer(ctx context.Context, obj *model.Offer) (*model.User, error) {
	if obj.DB.Consumer.ID == obj.DB.ConsumerID {
		return (&model.User{}).From(&obj.DB.Consumer)
	}

	user := db.User{}
	result := db.DB.Take(&user, "id = ?", obj.DB.ConsumerID)

	if result.Error != nil {
		return nil, fmt.Errorf("db take: %w", result.Error)
	}

	return (&model.User{}).From(&user)
}

func (r *offerResolver) Product(ctx context.Context, obj *model.Offer) (*model.Product, error) {
	if obj.DB.Product.ID == obj.DB.ProductID {
		return (&model.Product{}).From(&obj.DB.Product)
	}

	product := db.Product{}
	result := db.DB.Take(&product, "id = ?", obj.DB.ProductID)

	if result.Error != nil {
		return nil, fmt.Errorf("db take: %w", result.Error)
	}

	return (&model.Product{}).From(&product)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Offer returns generated.OfferResolver implementation.
func (r *Resolver) Offer() generated.OfferResolver { return &offerResolver{r} }

type mutationResolver struct{ *Resolver }
type offerResolver struct{ *Resolver }
