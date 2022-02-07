package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"auction-back/jwt"
	"context"
	"fmt"

	"github.com/teris-io/shortid"
)

func (r *mutationResolver) Register(ctx context.Context) (*model.RegisterResult, error) {
	id, err := shortid.Generate()

	if err != nil {
		return nil, err
	}

	user := db.User{
		ID: id,
	}

	result := db.DB.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	token, err := jwt.New(user.ID)

	if err != nil {
		return nil, err
	}

	return &model.RegisterResult{
		Token: token,
	}, nil
}

func (r *mutationResolver) IncreaseBalance(ctx context.Context, input model.IncreaseBalanceInput) (*model.IncreaseBalanceResult, error) {
	user := db.User{}

	result := db.DB.First(&user, "id = ?", input.UserID)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO: fix precision
	available := user.Available + input.Amount

	if available < 0 {
		return nil, fmt.Errorf("available balance cant be negative")
	}

	user.Available = available

	db.DB.Save(&user)

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

func (r *userResolver) Offers(ctx context.Context, obj *model.User) ([]*model.Offer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Products(ctx context.Context, obj *model.User, first *int, after *string) (*model.ProductConnection, error) {
	query := db.DB.Where("owner_id = ?", obj.ID).Order("id")

	fmt.Println("After query := ")

	if first != nil {
		if *first < 1 {
			return nil, fmt.Errorf("first must be positive")
		}
		query = query.Limit(*first + 1)
	}

	if after != nil {
		query.Where("id > ?", after)
	}

	var products []db.Product

	result := query.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(products) == 0 {
		return &model.ProductConnection{
			PageInfo: &model.PageInfo{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     nil,
				EndCursor:       nil,
			},
			Edges: make([]*model.ProductConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(products) > *first
		products = products[:len(products)-1]
	}

	edges := make([]*model.ProductConnectionEdge, 0, len(products))

	for _, product := range products {
		node, err := (&model.Product{}).From(&product)

		if err != nil {
			return nil, err
		}

		edges = append(edges, &model.ProductConnectionEdge{
			Cursor: product.ID,
			Node:   node,
		})
	}

	return &model.ProductConnection{
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &products[0].ID,
			EndCursor:       &products[len(products)-1].ID,
		},
		Edges: edges,
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
