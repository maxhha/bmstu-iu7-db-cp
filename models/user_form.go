package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
)

type UserForm struct {
	ID            string `json:"id"`
	UserID        string
	State         UserFormState
	Name          *string `json:"name"`
	Password      *string
	Phone         *string       `json:"phone"`
	Email         *string       `json:"email"`
	Currency      *CurrencyEnum `json:"currency"`
	DeclainReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}

func (f *UserForm) IsEditable() bool {
	return f.State == UserFormStateCreated || f.State == UserFormStateDeclained
}

func (f *UserFormFilled) From(form *UserForm) (*UserFormFilled, error) {
	var err error

	if form.Email == nil {
		err = multierror.Append(err, fmt.Errorf("no email"))
	} else {
		f.Email = *form.Email
	}

	if form.Phone == nil {
		err = multierror.Append(err, fmt.Errorf("no phone"))
	} else {
		f.Phone = *form.Phone
	}

	if form.Name == nil {
		err = multierror.Append(err, fmt.Errorf("no name"))
	} else {
		f.Name = *form.Name
	}

	if form.Currency == nil {
		err = multierror.Append(err, fmt.Errorf("no currency"))
	} else {
		f.Currency = *form.Currency
	}

	return f, err
}
