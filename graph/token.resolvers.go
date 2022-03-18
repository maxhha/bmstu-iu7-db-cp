package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/model"
	"context"
	"fmt"
	"time"
)

func (r *mutationResolver) CreateToken(ctx context.Context, input *model.CreateTokenInput) (*bool, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	action := db.TokenAction(input.Action.String())
	if err := r.Token.Validate(action, input.Data); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	token := db.Token{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)),
		Action:    action,
		Data:      input.Data,
		UserID:    viewer.ID,
	}

	if err := r.DB.Create(&token).Error; err != nil {
		return nil, err
	}

	if err := r.Token.Send(token); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	res := true
	return &res, nil
}
