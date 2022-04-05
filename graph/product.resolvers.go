package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input models.CreateProductInput) (*models.ProductResult, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if _, err := viewer.LastApprovedUserForm(r.DB); err != nil {
		return nil, fmt.Errorf("last approved user form: %w", err)
	}

	product := models.Product{
		Title:       input.Title,
		Description: input.Description,
		CreatorID:   viewer.ID,
	}

	if err := r.DB.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	product.Creator = *viewer

	return &models.ProductResult{
		Product: &product,
	}, nil
}

func (r *mutationResolver) OfferProduct(ctx context.Context, input models.OfferProductInput) (*models.OfferProductResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	// product := models.Product{}

	// result := r.DBTake(&product, "id = ?", input.ProductID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take: %w", result.Error)
	// }

	// if product.OwnerID != viewer.ID {
	// 	return nil, fmt.Errorf("viewer is not owner")
	// }

	// product.IsOnMarket = true

	// result = r.DBSave(&product)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db save: %w", result.Error)
	// }

	// p, err := (&models.Product{}).From(&product)

	// if err != nil {
	// 	return nil, fmt.Errorf("convert: %w", err)
	// }

	// go func() {
	// 	r.MarketLock.Lock()

	// 	for _, ch := range r.Market {
	// 		ch <- p
	// 	}

	// 	r.MarketLock.Unlock()
	// }()

	// return &models.OfferProductResult{Product: p}, nil
}

func (r *mutationResolver) TakeOffProduct(ctx context.Context, input models.TakeOffProductInput) (*models.TakeOffProductResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	// product := models.Product{}

	// result := r.DBTake(&product, "id = ?", input.ProductID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take: %w", result.Error)
	// }

	// if product.OwnerID != viewer.ID {
	// 	return nil, fmt.Errorf("viewer is not owner")
	// }

	// product.IsOnMarket = false

	// result = r.DBSave(&product)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db save: %w", result.Error)
	// }

	// p, err := (&models.Product{}).From(&product)

	// if err != nil {
	// 	return nil, fmt.Errorf("convert: %w", err)
	// }

	// return &models.TakeOffProductResult{Product: p}, nil
}

func (r *mutationResolver) SellProduct(ctx context.Context, input models.SellProductInput) (*models.SellProductResult, error) {
	panic(fmt.Errorf("not implemented"))

	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	// product := models.Product{}

	// if err := r.DBTake(&product, "id = ?", input.ProductID).Error; err != nil {
	// 	return nil, fmt.Errorf("db take: %w", err)
	// }

	// if product.OwnerID != viewer.ID {
	// 	return nil, fmt.Errorf("viewer is not owner")
	// }

	// offer := models.Offer{}

	// maxAmount := r.DBModel(&models.Offer{}).Select("max(amount)").Where("product_id = ?", product.ID)

	// if err := r.DB.Where("amount = (?)", maxAmount).First(&offer, "product_id = ?", product.ID).Error; err != nil {
	// 	return nil, fmt.Errorf("cant find max offer: %w", err)
	// }

	// err := r.DBTransaction(func(tx *gorm.DB) error {
	// 	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(viewer, "id = ?", viewer.ID).Error; err != nil {
	// 		return fmt.Errorf("lock viewer: %w", err)
	// 	}

	// 	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&product, "id = ?", product.ID).Error; err != nil {
	// 		return fmt.Errorf("lock product: %w", err)
	// 	}

	// 	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&offer, "id = ?", offer.ID).Error; err != nil {
	// 		return err
	// 	}

	// 	viewer.Available = viewer.Available + offer.Amount
	// 	product.OwnerID = offer.ConsumerID
	// 	product.IsOnMarket = false

	// 	if err := tx.Delete(&offer).Error; err != nil {
	// 		return fmt.Errorf("delete offer: %w", err)
	// 	}

	// 	if err := tx.Save(viewer).Error; err != nil {
	// 		return fmt.Errorf("save viewer: %w", err)
	// 	}

	// 	if err := tx.Save(&product).Error; err != nil {
	// 		return fmt.Errorf("save product: %w", err)
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	return nil, err
	// }

	// go func(id string) {
	// 	var offers []models.Offer
	// 	if err := r.DBFind(&offers, "product_id = ?", id).Error; err != nil {
	// 		fmt.Fprintf(gin.DefaultErrorWriter, "sell sideffect: db take: %v", err)
	// 		return
	// 	}

	// 	for _, offer := range offers {
	// 		o, err := (&models.Offer{}).From(&offer)

	// 		if err != nil {
	// 			fmt.Fprintf(gin.DefaultErrorWriter, "sell sideffect: convert: %v", err)
	// 			continue
	// 		}

	// 		if err := o.RemoveOffer(); err != nil {
	// 			fmt.Fprintf(gin.DefaultErrorWriter, "sell sideffect: remove offer: %v", err)
	// 			continue
	// 		}
	// 	}
	// }(product.ID)

	// p, err := (&models.Product{}).From(&product)

	// if err != nil {
	// 	return nil, err
	// }

	// return &models.SellProductResult{
	// 	Product: p,
	// }, nil
}

func (r *productResolver) Owner(ctx context.Context, obj *models.Product) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
	// if obj.DB.Owner.ID == obj.DB.OwnerID {
	// 	return (&models.User{}).From(&obj.DB.Owner)
	// }

	// owner := models.User{}
	// result := r.DBTake(&owner, "id = ?", obj.DB.OwnerID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take: %w", result.Error)
	// }

	// return (&models.User{}).From(&owner)
}

func (r *productResolver) TopOffer(ctx context.Context, obj *models.Product) (*models.Offer, error) {
	offer := models.Offer{}

	maxAmount := r.DB.Model(&models.Offer{}).Select("max(amount)").Where("product_id = ?", obj.ID)

	if err := r.DB.Where("amount = (?)", maxAmount).First(&offer, "product_id = ?", obj.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("cant find max offer: %w", err)
	}

	return &offer, nil
}

func (r *productResolver) Images(ctx context.Context, obj *models.Product) ([]*models.ProductImage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Offers(ctx context.Context, obj *models.Product, first *int, after *string) (*models.OffersConnection, error) {
	query := r.DB.Where("product_id = ?", obj.ID).Order("id")

	return OfferPagination(query, first, after)
}

func (r *queryResolver) MarketProducts(ctx context.Context, first *int, after *string) (*models.ProductsConnection, error) {
	query := r.DB.Where("is_on_market = true").Order("id")

	return ProductPagination(query, first, after)
}

func (r *subscriptionResolver) ProductOffered(ctx context.Context) (<-chan *models.Product, error) {
	ch := make(chan *models.Product, 1)

	r.MarketLock.Lock()
	chan_id := randString(6)
	r.Market[chan_id] = ch
	r.MarketLock.Unlock()

	go func() {
		<-ctx.Done()
		r.MarketLock.Lock()
		delete(r.Market, chan_id)
		r.MarketLock.Unlock()
	}()

	return ch, nil
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
