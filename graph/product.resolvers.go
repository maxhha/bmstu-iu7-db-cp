package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateProduct(ctx context.Context) (*models.ProductResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := r.DB.User().LastApprovedUserForm(viewer); err != nil {
		return nil, fmt.Errorf("last approved user form: %w", err)
	}

	product := models.Product{
		CreatorID: viewer.ID,
	}

	if err := r.DB.Product().Create(&product); err != nil {
		return nil, fmt.Errorf("db create: %w", err)
	}

	return &models.ProductResult{
		Product: &product,
	}, nil
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, input models.UpdateProductInput) (*models.ProductResult, error) {
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

	if !product.IsEditable() {
		return nil, fmt.Errorf("product is not editable")
	}

	product.Title = input.Title
	product.Description = input.Description

	if err := r.DB.Product().Update(&product); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.ProductResult{
		Product: &product,
	}, nil
}

func (r *mutationResolver) RequestModerateProduct(ctx context.Context, input models.ProductInput) (bool, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return false, err
	}

	product, err := r.DB.Product().Get(input.ProductID)
	if err != nil {
		return false, fmt.Errorf("db get: %w", err)
	}

	if err := isProductOwner(r.DB, viewer, product); err != nil {
		return false, err
	}

	if !product.IsEditable() {
		return false, fmt.Errorf("product is not editable")
	}

	data := map[string]interface{}{"productID": product.ID}
	if err := r.TokenPort.Create(models.TokenActionModerateProduct, viewer, data); err != nil {
		return false, fmt.Errorf("token create: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) ApproveModerateProduct(ctx context.Context, input models.TokenInput) (*models.ProductResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	token, err := r.TokenPort.Activate(models.TokenActionModerateProduct, input.Token, viewer)
	if err != nil {
		return nil, fmt.Errorf("token activate: %w", err)
	}

	productID, ok := token.Data["productID"].(string)
	if !ok {
		return nil, fmt.Errorf("no productID in token")
	}

	product, err := r.DB.Product().Get(productID)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	if err := isProductOwner(r.DB, viewer, product); err != nil {
		return nil, err
	}

	if !product.IsEditable() {
		return nil, fmt.Errorf("product is not editable")
	}

	product.State = models.ProductStateModerating

	if err := r.DB.Product().Update(&product); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.ProductResult{
		Product: &product,
	}, nil
}

func (r *mutationResolver) ApproveProduct(ctx context.Context, input models.ProductInput) (*models.ProductResult, error) {
	product, err := r.DB.Product().Get(input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	if product.State != models.ProductStateModerating {
		return nil, fmt.Errorf("state is not %s", models.ProductStateModerating)
	}

	product.State = models.ProductStateApproved

	if err := r.DB.Product().Update(&product); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.ProductResult{
		Product: &product,
	}, nil
}

func (r *mutationResolver) DeclainProduct(ctx context.Context, input models.DeclineProductInput) (*models.ProductResult, error) {
	product, err := r.DB.Product().Get(input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	if product.State != models.ProductStateModerating && product.State != models.ProductStateApproved {
		return nil, fmt.Errorf("state is not %s", models.ProductStateModerating)
	}

	product.State = models.ProductStateDeclained
	product.DeclainReason = input.DeclainReason

	if err := r.DB.Product().Update(&product); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.ProductResult{
		Product: &product,
	}, nil
}

func (r *mutationResolver) OfferProduct(ctx context.Context, input models.ProductInput) (*models.OfferProductResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, ErrUnauthorized
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

func (r *mutationResolver) TakeOffProduct(ctx context.Context, input models.ProductInput) (*models.TakeOffProductResult, error) {
	panic(fmt.Errorf("not implemented"))
	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, ErrUnauthorized
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

func (r *mutationResolver) SellProduct(ctx context.Context, input models.ProductInput) (*models.SellProductResult, error) {
	panic(fmt.Errorf("not implemented"))

	// viewer := auth.ForViewer(ctx)

	// if viewer == nil {
	// 	return nil, ErrUnauthorized
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
	if obj == nil {
		return nil, fmt.Errorf("product is nil")
	}

	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, fmt.Errorf("for viewer: %w", err)
	}

	if err := r.isProductOwnerOrManager(viewer, *obj); err != nil {
		return nil, fmt.Errorf("owner or manager: %w", err)
	}

	owner, err := r.DB.Product().GetOwner(*obj)
	if err != nil {
		return nil, fmt.Errorf("get owner: %w", err)
	}

	return &owner, nil
}

func (r *productResolver) Creator(ctx context.Context, obj *models.Product) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) TopOffer(ctx context.Context, obj *models.Product) (*models.Offer, error) {
	panic(fmt.Errorf("not implemented"))
	// offer := models.Offer{}

	// maxAmount := r.DB.Model(&models.Offer{}).Select("max(amount)").Where("product_id = ?", obj.ID)

	// if err := r.DB.Where("amount = (?)", maxAmount).First(&offer, "product_id = ?", obj.ID).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, nil
	// 	}

	// 	return nil, fmt.Errorf("cant find max offer: %w", err)
	// }

	// return &offer, nil
}

func (r *productResolver) Images(ctx context.Context, obj *models.Product) ([]*models.ProductImage, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Offers(ctx context.Context, obj *models.Product, first *int, after *string) (*models.OffersConnection, error) {
	panic(fmt.Errorf("not implemented"))
	// query := r.DB.Where("product_id = ?", obj.ID).Order("id")

	// return OfferPagination(query, first, after)
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string, filter *models.ProductsFilter) (*models.ProductsConnection, error) {
	connection, err := r.DB.Product().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}

	return &connection, nil
}

func (r *queryResolver) MarketProducts(ctx context.Context, first *int, after *string) (*models.ProductsConnection, error) {
	panic(fmt.Errorf("not implemented"))
	// query := r.DB.Where("is_on_market = true").Order("id")

	// return ProductPagination(query, first, after)
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

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type productResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
