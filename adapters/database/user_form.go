package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: check if gorm.DB.Update updates objects field UpdatedAt
//go:generate go run ../../codegen/gormdbops/main.go --out user_form_gen.go --model UserForm --methods Get,Pagination,Update,Create

type UserForm struct {
	ID            string `gorm:"default:generated();"`
	UserID        string
	State         models.UserFormState `gorm:"default:'CREATED';"`
	Name          *string
	Password      *string
	Phone         *string
	Email         *string
	Currency      *models.CurrencyEnum
	DeclainReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (u *UserForm) into() models.UserForm {
	return models.UserForm{
		ID:            u.ID,
		UserID:        u.UserID,
		State:         u.State,
		Name:          u.Name,
		Password:      u.Password,
		Phone:         u.Phone,
		Email:         u.Email,
		Currency:      u.Currency,
		DeclainReason: u.DeclainReason,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		DeletedAt:     sql.NullTime(u.DeletedAt),
	}
}

func (f *UserForm) copy(form *models.UserForm) {
	if form == nil {
		return
	}
	f.ID = form.ID
	f.UserID = form.UserID
	f.State = form.State
	f.Name = form.Name
	f.Password = form.Password
	f.Phone = form.Phone
	f.Email = form.Email
	f.Currency = form.Currency
	f.DeclainReason = form.DeclainReason
	f.CreatedAt = form.CreatedAt
	f.UpdatedAt = form.UpdatedAt
	f.DeletedAt = gorm.DeletedAt(form.DeletedAt)
}

var userFormFieldToColumn = map[ports.UserFormField]string{
	ports.UserFormFieldCreatedAt: "created_at",
}

func (d *userFormDB) filter(query *gorm.DB, config *models.UserFormsFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.ID) > 0 {
		query = query.Where("id IN ?", config.ID)
	}

	if len(config.UserID) > 0 {
		query = query.Where("user_id IN ?", config.UserID)
	}

	if len(config.State) > 0 {
		query = query.Where("state IN ?", config.State)
	}

	return query
}

func (d *userFormDB) Take(config ports.UserFormTakeConfig) (models.UserForm, error) {
	query := d.filter(d.db, &config.UserFormsFilter)

	if config.OrderBy != "" {
		column, ok := userFormFieldToColumn[config.OrderBy]
		if !ok {
			return models.UserForm{}, fmt.Errorf("unknown field '%s'", config.OrderBy)
		}

		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   config.OrderDesc,
		})
	}

	userForm := UserForm{}
	if err := query.Take(&userForm).Error; err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return userForm.into(), nil
}

func approvedOrFirstUserFormFilter(query *gorm.DB) *gorm.DB {
	return query.
		Where(`(
			state = 'APPROVED'
			OR (SELECT COUNT(1) FROM user_forms u WHERE "user_forms"."user_id" = u.user_id) = 1
		)`)
}

func (d *userFormDB) GetLoginForm(input models.LoginInput) (models.UserForm, error) {
	form := UserForm{}
	query := approvedOrFirstUserFormFilter(d.db)
	err := query.
		Where(
			"name = @username OR email = @username OR phone = @username",
			sql.Named("username", input.Username),
		).
		Where(
			"password IS NOT NULL",
		).
		Take(
			&form,
		).Error

	if err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return form.into(), nil
}
