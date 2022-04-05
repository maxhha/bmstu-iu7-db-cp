package models

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

type UserForm struct {
	gorm.Model
	ID            string `json:"id" gorm:"default:generated();"`
	UserID        string
	User          User
	State         UserFormState `gorm:"default:'CREATED';"`
	Name          *string       `json:"name"`
	Password      *string
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	DeclainReason *string
}

func (f *UserForm) MostRelevantFilter(db *gorm.DB) *gorm.DB {
	return db.Model(f).
		Where(`(
			state = 'APPROVED' 
			OR (SELECT COUNT(1) FROM user_forms u WHERE "user_forms"."user_id" = u.user_id) = 1
		)`)
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

	return f, err
}
