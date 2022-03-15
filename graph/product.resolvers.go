package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.CreateProductResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	// id, err := shortid.Generate()

	// if err != nil {
	// 	return nil, fmt.Errorf("shortid: %w", err)
	// }

	// product := db.Product{
	// 	ID:          id,
	// 	Name:        input.Name,
	// 	Description: input.Description,
	// 	OwnerID:     viewer.ID,
	// 	Owner:       *viewer,
	// }

	// result := db.DB.Create(&product)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db create: %w", result.Error)
	// }

	// p, err := (&model.Product{}).From(&product)

	// if err != nil {
	// 	return nil, fmt.Errorf("convert: %w", result.Error)
	// }

	// return &model.CreateProductResult{
	// 	Product: p,
	// }, nil
}

func (r *mutationResolver) OfferProduct(ctx context.Context, input model.OfferProductInput) (*model.OfferProductResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	// product := db.Product{}

	// result := db.DB.Take(&product, "id = ?", input.ProductID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take: %w", result.Error)
	// }

	// if product.OwnerID != viewer.ID {
	// 	return nil, fmt.Errorf("viewer is not owner")
	// }

	// product.IsOnMarket = true

	// result = db.DB.Save(&product)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db save: %w", result.Error)
	// }

	// p, err := (&model.Product{}).From(&product)

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

	// return &model.OfferProductResult{Product: p}, nil
}

func (r *mutationResolver) TakeOffProduct(ctx context.Context, input model.TakeOffProductInput) (*model.TakeOffProductResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	// product := db.Product{}

	// result := db.DB.Take(&product, "id = ?", input.ProductID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take: %w", result.Error)
	// }

	// if product.OwnerID != viewer.ID {
	// 	return nil, fmt.Errorf("viewer is not owner")
	// }

	// product.IsOnMarket = false

	// result = db.DB.Save(&product)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db save: %w", result.Error)
	// }

	// p, err := (&model.Product{}).From(&product)

	// if err != nil {
	// 	return nil, fmt.Errorf("convert: %w", err)
	// }

	// return &model.TakeOffProductResult{Product: p}, nil
}

func (r *mutationResolver) SellProduct(ctx context.Context, input model.SellProductInput) (*model.SellProductResult, error) {
	panic(fmt.Errorf("not implemented"))

	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	// product := db.Product{}

	// if err := db.DB.Take(&product, "id = ?", input.ProductID).Error; err != nil {
	// 	return nil, fmt.Errorf("db take: %w", err)
	// }

	// if product.OwnerID != viewer.ID {
	// 	return nil, fmt.Errorf("viewer is not owner")
	// }

	// offer := db.Offer{}

	// maxAmount := db.DB.Model(&db.Offer{}).Select("max(amount)").Where("product_id = ?", product.ID)

	// if err := db.DB.Where("amount = (?)", maxAmount).First(&offer, "product_id = ?", product.ID).Error; err != nil {
	// 	return nil, fmt.Errorf("cant find max offer: %w", err)
	// }

	// err := db.DB.Transaction(func(tx *gorm.DB) error {
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
	// 	var offers []db.Offer
	// 	if err := db.DB.Find(&offers, "product_id = ?", id).Error; err != nil {
	// 		fmt.Fprintf(gin.DefaultErrorWriter, "sell sideffect: db take: %v", err)
	// 		return
	// 	}

	// 	for _, offer := range offers {
	// 		o, err := (&model.Offer{}).From(&offer)

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

	// p, err := (&model.Product{}).From(&product)

	// if err != nil {
	// 	return nil, err
	// }

	// return &model.SellProductResult{
	// 	Product: p,
	// }, nil
}

func (r *productResolver) Title(ctx context.Context, obj *model.Product) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Owner(ctx context.Context, obj *model.Product) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
	// if obj.DB.Owner.ID == obj.DB.OwnerID {
	// 	return (&model.User{}).From(&obj.DB.Owner)
	// }

	// owner := db.User{}
	// result := db.DB.Take(&owner, "id = ?", obj.DB.OwnerID)

	// if result.Error != nil {
	// 	return nil, fmt.Errorf("db take: %w", result.Error)
	// }

	// return (&model.User{}).From(&owner)
}

func (r *productResolver) Creator(ctx context.Context, obj *model.Product) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) TopOffer(ctx context.Context, obj *model.Product) (*model.Offer, error) {
	offer := db.Offer{}

	maxAmount := db.DB.Model(&db.Offer{}).Select("max(amount)").Where("product_id = ?", obj.ID)

	if err := db.DB.Where("amount = (?)", maxAmount).First(&offer, "product_id = ?", obj.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("cant find max offer: %w", err)
	}

	return (&model.Offer{}).From(&offer)
}

func (r *productResolver) Images(ctx context.Context, obj *model.Product) ([]*model.ProductImage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Offers(ctx context.Context, obj *model.Product, first *int, after *string) (*model.OffersConnection, error) {
	query := db.DB.Where("product_id = ?", obj.ID).Order("id")

	return OfferPagination(query, first, after)
}

func (r *queryResolver) MarketProducts(ctx context.Context, first *int, after *string) (*model.ProductsConnection, error) {
	query := db.DB.Where("is_on_market = true").Order("id")

	return ProductPagination(query, first, after)
}

func (r *subscriptionResolver) ProductOffered(ctx context.Context) (<-chan *model.Product, error) {
	ch := make(chan *model.Product, 1)

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
