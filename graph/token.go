package graph

import (
	"auction-back/auth"
	"auction-back/db"
	"context"
	"fmt"
)

var validateTokenData = map[db.TokenAction]func(data map[string]interface{}) error{
	db.TokenActionApproveUserEmail: func(data map[string]interface{}) error {
		email, found := data["email"]

		if !found {
			return fmt.Errorf("no email in data")
		}

		_, ok := email.(string)
		if !ok {
			return fmt.Errorf("email in data is not string")
		}

		return nil
	},
	db.TokenActionApproveUserPhone: func(data map[string]interface{}) error {
		phone, found := data["phone"]

		if !found {
			return fmt.Errorf("no phone in data")
		}

		_, ok := phone.(string)
		if !ok {
			return fmt.Errorf("phone in data is not string")
		}

		return nil
	},
}

func getTokenCreator(ctx context.Context) (*db.TokenCreator, error) {
	user := auth.ForViewer(ctx)
	guest := auth.ForGuest(ctx)

	if user == nil && guest == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	var id *string
	var condition string

	if user != nil {
		id = &user.ID
		condition = "user_id = ? AND guest_id IS NULL"
	} else {
		id = &guest.ID
		condition = "guest_id = ? AND user_id IS NULL"
	}

	creator := db.TokenCreator{}

	if err := db.DB.Take(&creator, condition, id).Error; err != nil {
		return nil, fmt.Errorf("take creator: %w", err)
	}

	return &creator, nil
}
