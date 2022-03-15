package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/db"
	"auction-back/graph/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (r *mutationResolver) CreateToken(ctx context.Context, input *model.CreateTokenInput) (*bool, error) {
	creator, err := getTokenCreator(ctx)

	if err != nil {
		return nil, err
	}

	tokenAction := db.TokenAction(input.Action.String())
	validate, found := validateTokenData[tokenAction]
	if !found {
		return nil, fmt.Errorf("not found validator for action")
	}

	if err := validate(input.Data); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	token := db.Token{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)),
		Action:    tokenAction,
		Data:      input.Data,
		CreatorID: creator.ID,
	}

	if err := db.DB.Create(&token).Error; err != nil {
		return nil, err
	}

	// TODO send token somehow
	fmt.Println("token:", token.ID)

	res := true

	return &res, nil
}

func (r *mutationResolver) ActivateToken(ctx context.Context, input *model.ActivateTokenInput) (*bool, error) {
	creator, err := getTokenCreator(ctx)

	if err != nil {
		return nil, err
	}

	token := db.Token{}

	if err := db.DB.Take(&token, "id = ?", input.Token).Error; err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	if token.CreatorID != creator.ID {
		return nil, fmt.Errorf("creator is other")
	}

	token.ActivatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err := db.DB.Save(&token).Error; err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	res := true

	return &res, nil
}
