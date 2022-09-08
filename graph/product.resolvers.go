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

	if err := IsProductOwner(r.DB, viewer, product); err != nil {
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

	if err := IsProductOwner(r.DB, viewer, product); err != nil {
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

	if err := IsProductOwner(r.DB, viewer, product); err != nil {
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
	user, err := r.DB.User().Get(obj.CreatorID)
	if err != nil {
		return nil, fmt.Errorf("r.DB.User().Get: %w", err)
	}

	return &user, nil
}

func (r *productResolver) Images(ctx context.Context, obj *models.Product) ([]*models.ProductImage, error) {
	return nil, nil
}

func (r *productResolver) Auctions(ctx context.Context, obj *models.Product, first *int, after *string, filter *models.AuctionsFilter) (*models.AuctionsConnection, error) {
	if filter == nil {
		filter = &models.AuctionsFilter{}
	}

	filter.ProductIDs = []string{obj.ID}
	conn, err := r.DB.Auction().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Auction().Pagination: %w", err)
	}

	return &conn, nil
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string, filter *models.ProductsFilter) (*models.ProductsConnection, error) {
	connection, err := r.DB.Product().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}

	return &connection, nil
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }
