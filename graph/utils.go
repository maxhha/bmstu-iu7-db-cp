package graph

import (
	"auction-back/models"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func moneyMapToArray(currencyMap map[models.CurrencyEnum]models.Money) []*models.Money {
	result := make([]*models.Money, 0, len(currencyMap))

	for _, m := range currencyMap {
		tmp := m
		result = append(result, &tmp)
	}

	return result
}

// func (r *balanceResolver) Blocked(ctx context.Context, obj *models.Balance) (float64, error) {
// 	var blocked *float64
// 	if err := r.DBModel(&db.Offer{}).Select("sum(amount)").Where("consumer_id = ?", obj.DB.ID).Scan(&blocked).Error; err != nil {
// 		return 0, err
// 	}

// 	if blocked == nil {
// 		return 0, nil
// 	}

// 	return *blocked, nil
// }
// func (r *userResolver) Balance(ctx context.Context, obj *models.User) (*models.Balance, error) {
// 	viewer := auth.ForViewer(ctx)

// 	if viewer == nil {
// 		return nil, ErrUnauthorized
// 	}

// 	if viewer.ID != obj.ID {
// 		return nil, fmt.Errorf("denied")
// 	}

// 	return (&models.Balance{}).From(obj.DB)
// }
// func (r *Resolver) Balance() generated.BalanceResolver { return &balanceResolver{r} }

// type balanceResolver struct{ *Resolver }
