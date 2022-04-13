package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
	"time"
)

func (r *mutationResolver) CreateOffer(ctx context.Context, input models.CreateOfferInput) (*models.CreateOfferResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, ErrUnauthorized
	// }

	// var product models.Product

	// result := r.DBTake(&product, "id = ?", input.ProductID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take product: %w", result.Error)
	// }

	// if !product.IsOnMarket {
	// 	return nil, fmt.Errorf("product is not on market")
	// }

	// id, err := shortid.Generate()

	// if err != nil {
	// 	return nil, fmt.Errorf("shortid: %w", err)
	// }

	// offer := db.Offer{
	// 	ID:         id,
	// 	Amount:     input.Amount,
	// 	ProductID:  input.ProductID,
	// 	Product:    product,
	// 	ConsumerID: viewer.ID,
	// 	Consumer:   *viewer,
	// }

	// err = r.DBTransaction(func(tx *gorm.DB) error {
	// 	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(viewer, "id = ?", viewer.ID).Error; err != nil {
	// 		return err
	// 	}

	// 	if viewer.Available < input.Amount {
	// 		return fmt.Errorf("not enough amount")
	// 	}

	// 	// FIXME: presicion
	// 	viewer.Available = viewer.Available - input.Amount

	// 	if err := tx.Save(&viewer).Error; err != nil {
	// 		return err
	// 	}

	// 	if err := tx.Save(&offer).Error; err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	return nil, fmt.Errorf("db save: %w", err)
	// }

	// o, err := (&models.Offer{}).From(&offer)

	// if err != nil {
	// 	return nil, err
	// }

	// return &models.CreateOfferResult{
	// 	Offer: o,
	// }, nil
}

func (r *mutationResolver) RemoveOffer(ctx context.Context, input models.RemoveOfferInput) (*models.RemoveOfferResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, ErrUnauthorized
	// }

	// offer := db.Offer{}

	// if err := r.DBTake(&offer, "id = ?", input.OfferID).Error; err != nil {
	// 	return nil, fmt.Errorf("db take: %w", err)
	// }

	// if offer.ConsumerID != viewer.ID {
	// 	return nil, fmt.Errorf("denied")
	// }

	// o, err := (&models.Offer{}).From(&offer)

	// if err != nil {
	// 	return nil, fmt.Errorf("convert: %w", err)
	// }

	// if err := o.RemoveOffer(); err != nil {
	// 	return nil, fmt.Errorf("offer remove: %w", err)
	// }

	// return &models.RemoveOfferResult{
	// 	Status: "success",
	// }, nil
}

func (r *offerResolver) State(ctx context.Context, obj *models.Offer) (models.OfferStateEnum, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) FailReason(ctx context.Context, obj *models.Offer) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) Moneys(ctx context.Context, obj *models.Offer) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) Transactions(ctx context.Context, obj *models.Offer) ([]*models.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

// Offer returns generated.OfferResolver implementation.
func (r *Resolver) Offer() generated.OfferResolver { return &offerResolver{r} }

type offerResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *offerResolver) User(ctx context.Context, obj *models.Offer) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *offerResolver) Product(ctx context.Context, obj *models.Offer) (*models.Product, error) {
	panic(fmt.Errorf("not implemented"))
	// if obj.DB.Product.ID == obj.DB.ProductID {
	// 	return (&models.Product{}).From(&obj.DB.Product)
	// }

	// product := models.Product{}
	// result := r.DBTake(&product, "id = ?", obj.DB.ProductID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take: %w", result.Error)
	// }

	// return (&models.Product{}).From(&product)
}
func (r *offerResolver) CreatedAt(ctx context.Context, obj *models.Offer) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *offerResolver) DeleteOnSell(ctx context.Context, obj *models.Offer) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}
