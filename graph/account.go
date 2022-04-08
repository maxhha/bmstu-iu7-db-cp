package graph

import (
	"auction-back/models"
	"context"
	"fmt"
)

func (r *accountResolver) Bank(ctx context.Context, obj *models.Account) (*models.Bank, error) {
	if obj == nil {
		return nil, fmt.Errorf("account is nil")
	}

	bank, err := r.DB.Bank().Get(obj.BankID)
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (r *accountResolver) Transactions(ctx context.Context, obj *models.Account, first *int, after *string) (*models.TransactionsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Account() *accountResolver { return &accountResolver{r} }

type accountResolver struct{ *Resolver }
