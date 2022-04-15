package models

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shopspring/decimal"
)

const dateTimeLayout = `2006-01-02T15:04:05.000Z07:00`

func MarshalDateTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(`"`))
		w.Write([]byte(t.UTC().Format(dateTimeLayout)))
		w.Write([]byte(`"`))
	})
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {
	str, ok := v.(string)

	if !ok {
		return time.Time{}, fmt.Errorf("convert to string")
	}

	t, err := time.Parse(dateTimeLayout, str)

	if err != nil {
		return time.Time{}, fmt.Errorf("time parse: %w", err)
	}

	return t, nil
}

type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency CurrencyEnum    `json:"currency"`
}

func (m *MoneyInput) IntoPtr() *Money {
	if m == nil {
		return nil
	}

	return &Money{
		Amount:   m.Amount,
		Currency: m.Currency,
	}
}

func MarshalDecimal(d decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(d.String()))
	})
}

func UnmarshalDecimal(v interface{}) (decimal.Decimal, error) {
	valf, okf := v.(float64)
	if okf {
		return decimal.NewFromFloatWithExponent(valf, -2), nil
	}

	vali, oki := v.(int64)
	if oki {
		return decimal.NewFromInt(vali), nil
	}

	return decimal.Decimal{}, fmt.Errorf("fail convert to float or int %#v", v)
}
