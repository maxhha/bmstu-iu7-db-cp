package models

// type Offer struct {
// 	ID        string  `json:"id"`
// 	Amount    float64 `json:"amount"`
// 	CreatedAt string  `json:"createdAt"`
// 	DB        *db.Offer
// }

// func (o *Offer) From(offer *db.Offer) (*Offer, error) {
// 	o.ID = offer.ID
// 	// o.Amount = offer.Amount
// 	o.CreatedAt = offer.CreatedAt.String()
// 	o.DB = offer

// 	return o, nil
// }

// func (offer *Offer) RemoveOffer() error {
// 	err := r.DBTransaction(func(tx *gorm.DB) error {
// 		consumer := models.User{}

// 		// if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&consumer, "id = ?", offer.DB.ConsumerID).Error; err != nil {
// 		// return fmt.Errorf("lock viewer: %w", err)
// 		// }

// 		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(offer.DB, "id = ?", offer.ID).Error; err != nil {
// 			return fmt.Errorf("lock offer: %w", err)
// 		}

// 		// TODO: fix precision
// 		// consumer.Available = consumer.Available + offer.DB.Amount

// 		if err := tx.Delete(offer.DB).Error; err != nil {
// 			return fmt.Errorf("db delete: %w", err)
// 		}

// 		if err := tx.Save(consumer).Error; err != nil {
// 			return fmt.Errorf("db save: %w", err)
// 		}

// 		return nil
// 	})

// 	return err
// }
