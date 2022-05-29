package main

import (
	"auction-back/models"
	"bytes"
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type Time time.Time

func (t *Time) UnmarshalJSON(bs []byte) error {
	var v interface{}

	if err := json.Unmarshal(bs, &v); err != nil {
		return err
	}

	if v == nil {
		return nil
	}

	tm, err := models.UnmarshalTime(v)
	if err != nil {
		return err
	}
	*t = Time(tm)

	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	models.MarshalTime(time.Time(t)).MarshalGQL(&buf)
	return buf.Bytes(), nil
}

func (t *Time) Into() *time.Time {
	if t == nil {
		return nil
	}
	tm := time.Time(*t)
	return &tm
}

func (obj *Time) From(ent *time.Time) *Time {
	if ent == nil {
		return nil
	}

	*obj = Time(*ent)
	return obj
}

type Transaction struct {
	ID            int                     `json:"id"`
	Date          *Time                   `json:"date"`
	State         models.TransactionState `json:"state"`
	Type          models.TransactionType  `json:"type"`
	Currency      models.CurrencyEnum     `json:"currency"`
	Amount        decimal.Decimal         `json:"amount"`
	Error         *string                 `json:"error"`
	AccountFromID *string                 `json:"accountFrom"`
	AccountToID   *string                 `json:"accountTo"`
	OfferID       *string                 `json:"offer"`
	CreatedAt     Time                    `json:"createdAt"`
	UpdatedAt     Time                    `json:"updatedAt"`
	DeletedAt     *Time                   `json:"deletedAt"`
}

func (obj *Transaction) From(ent *models.Transaction) *Transaction {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID
	obj.Date = (&Time{}).From(ent.Date)
	obj.State = ent.State
	obj.Type = ent.Type
	obj.Currency = ent.Currency
	obj.Amount = ent.Amount
	obj.Error = ent.Error
	obj.AccountFromID = ent.AccountFromID
	obj.AccountToID = ent.AccountToID
	obj.OfferID = ent.OfferID
	obj.CreatedAt.From(&ent.CreatedAt)
	obj.UpdatedAt.From(&ent.UpdatedAt)
	obj.DeletedAt = (&Time{}).From(ent.DeletedAt)

	return obj
}
