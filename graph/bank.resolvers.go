package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/models"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateBank(ctx context.Context, input models.CreateBankInput) (*models.BankResult, error) {
	bank := models.Bank{
		Name:                 input.Name,
		Bic:                  input.Bic,
		CorrespondentAccount: input.CorrespondentAccount,
		Inn:                  input.Inn,
		Kpp:                  input.Kpp,
	}

	if err := r.DB.Bank().Create(&bank); err != nil {
		return nil, fmt.Errorf("db bank create: %w", err)
	}

	return &models.BankResult{
		Bank: &bank,
	}, nil
}

func (r *mutationResolver) UpdateBank(ctx context.Context, input models.UpdateBankInput) (*models.BankResult, error) {
	bank, err := r.DB.Bank().Get(input.BankID)
	if err != nil {
		return nil, fmt.Errorf("db bank get: %w", err)
	}

	bank.Name = input.Name
	bank.Bic = input.Bic
	bank.CorrespondentAccount = input.CorrespondentAccount
	bank.Inn = input.Inn
	bank.Kpp = input.Kpp

	if err := r.DB.Bank().Update(&bank); err != nil {
		return nil, fmt.Errorf("db bank update: %w", err)
	}

	return &models.BankResult{
		Bank: &bank,
	}, nil
}

func (r *queryResolver) Banks(ctx context.Context, first *int, after *string, filter *models.BanksFilter) (*models.BanksConnection, error) {
	conn, err := r.DB.Bank().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db bank pagination: %w", err)
	}

	return &conn, nil
}
