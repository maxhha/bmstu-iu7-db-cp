package model

import (
	"auction-back/db"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func (f *UserFormFilled) From(form *db.UserForm) (*UserFormFilled, error) {
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

	return f, err
}
