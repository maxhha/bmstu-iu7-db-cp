package graph

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// func (r *balanceResolver) Blocked(ctx context.Context, obj *model.Balance) (float64, error) {
// 	var blocked *float64
// 	if err := db.DB.Model(&db.Offer{}).Select("sum(amount)").Where("consumer_id = ?", obj.DB.ID).Scan(&blocked).Error; err != nil {
// 		return 0, err
// 	}

// 	if blocked == nil {
// 		return 0, nil
// 	}

// 	return *blocked, nil
// }
// func (r *userResolver) Balance(ctx context.Context, obj *model.User) (*model.Balance, error) {
// 	viewer := auth.ForViewer(ctx)

// 	if viewer == nil {
// 		return nil, fmt.Errorf("unauthorized")
// 	}

// 	if viewer.ID != obj.ID {
// 		return nil, fmt.Errorf("denied")
// 	}

// 	return (&model.Balance{}).From(obj.DB)
// }
// func (r *Resolver) Balance() generated.BalanceResolver { return &balanceResolver{r} }

// type balanceResolver struct{ *Resolver }
