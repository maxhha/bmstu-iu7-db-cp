package models

import "time"

type DealState struct {
	ID        string        `json:"id"`
	State     DealStateEnum `json:"state"`
	CreatorID *string
	OfferID   string
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
