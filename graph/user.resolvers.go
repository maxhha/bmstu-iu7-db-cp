package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/model"
	"auction-back/jwt"
	"context"
	"fmt"
	"time"
)

func (r *mutationResolver) CreateUser(ctx context.Context) (*model.CreateUserResult, error) {
	user := db.User{}

	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	token, err := jwt.NewUser(user.ID)

	if err != nil {
		return nil, err
	}

	return &model.CreateUserResult{
		Token: token,
	}, nil
}

func (r *mutationResolver) IncreaseBalance(ctx context.Context, input model.IncreaseBalanceInput) (*model.IncreaseBalanceResult, error) {
	user := db.User{}

	result := r.DB.First(&user, "id = ?", input.UserID)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO: fix precision
	// available := user.Available.input.Amount

	// if available < 0 {
	// 	return nil, fmt.Errorf("available balance cant be negative")
	// }

	// user.Available = available

	r.DB.Save(&user)

	u, err := (&model.User{}).From(&user)

	if err != nil {
		return nil, err
	}

	return &model.IncreaseBalanceResult{
		User: u,
	}, nil
}

func (r *queryResolver) Viewer(ctx context.Context) (*model.User, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, nil
	}

	return (&model.User{}).From(viewer)
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.

func (r *userResolver) Email(ctx context.Context, obj *model.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) Phone(ctx context.Context, obj *model.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) Name(ctx context.Context, obj *model.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) BlockedUntil(ctx context.Context, obj *model.User) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) IsAdmin(ctx context.Context, obj *model.User) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) Available(ctx context.Context, obj *model.User) ([]*model.Money, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) Blocked(ctx context.Context, obj *model.User) ([]*model.Money, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) Accounts(ctx context.Context, obj *model.User, first *int, after *string) (*model.UserAccountsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *userResolver) Offers(ctx context.Context, obj *model.User, first *int, after *string) (*model.OffersConnection, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if viewer.ID != obj.ID {
		return nil, fmt.Errorf("denied")
	}

	query := db.DB.Where("consumer_id = ?", obj.ID).Order("id")

	return OfferPagination(query, first, after)
}
func (r *userResolver) Products(ctx context.Context, obj *model.User, first *int, after *string) (*model.ProductsConnection, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if viewer.ID != obj.ID {
		return nil, fmt.Errorf("denied")
	}

	query := db.DB.Where("owner_id = ?", obj.ID).Order("id")

	return ProductPagination(query, first, after)
}

type userResolver struct{ *Resolver }
