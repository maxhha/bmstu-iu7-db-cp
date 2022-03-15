package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"auction-back/jwt"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (r *mutationResolver) RegisterGuest(ctx context.Context) (*model.RegisterGuestResult, error) {
	guest := db.Guest{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(6)),
	}

	fmt.Printf("query: %s\n", db.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(&guest)
	}))

	if err := db.DB.Create(&guest).Error; err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	fmt.Printf("guest_id = %v\n", guest.ID)

	token, err := jwt.NewGuest(guest.ID, guest.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return &model.RegisterGuestResult{
		Token: token,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
