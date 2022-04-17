package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/graph/generated"
	"auction-back/models"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateNominalAccount(ctx context.Context, input models.CreateNominalAccountInput) (*models.NominalAccountResult, error) {
	account := models.NominalAccount{
		Name:          input.Name,
		Receiver:      input.Receiver,
		AccountNumber: input.AccountNumber,
		BankID:        input.BankID,
	}

	if err := r.DB.NominalAccount().Create(&account); err != nil {
		return nil, fmt.Errorf("db nominal account create: %w", err)
	}

	return &models.NominalAccountResult{
		NominalAccount: &account,
	}, nil
}

func (r *mutationResolver) UpdateNominalAccount(ctx context.Context, input models.UpdateNominalAccountInput) (*models.NominalAccountResult, error) {
	account, err := r.DB.NominalAccount().Get(input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("db nominal account get: %w", err)
	}

	account.Name = input.Name
	account.Receiver = input.Receiver
	account.AccountNumber = input.AccountNumber
	account.BankID = input.BankID

	if err := r.DB.NominalAccount().Update(&account); err != nil {
		return nil, fmt.Errorf("db nominal account create: %w", err)
	}

	return &models.NominalAccountResult{
		NominalAccount: &account,
	}, nil
}

func (r *nominalAccountResolver) Bank(ctx context.Context, obj *models.NominalAccount) (*models.Bank, error) {
	if obj == nil {
		return nil, nil
	}

	bank, err := r.DB.Bank().Get(obj.BankID)
	if err != nil {
		return nil, fmt.Errorf("db bank get: %w", err)
	}

	return &bank, nil
}

func (r *queryResolver) NominalAccounts(ctx context.Context, first *int, after *string, filter *models.NominalAccountsFilter) (*models.NominalAccountsConnection, error) {
	conn, err := r.DB.NominalAccount().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db nominal account pagination: %w", err)
	}

	return &conn, nil
}

// NominalAccount returns generated.NominalAccountResolver implementation.
func (r *Resolver) NominalAccount() generated.NominalAccountResolver {
	return &nominalAccountResolver{r}
}

type nominalAccountResolver struct{ *Resolver }
