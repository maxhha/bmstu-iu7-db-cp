package database

import (
	"auction-back/models"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

//go:generate go run ../../codegen/gormdbops/main.go --out user_gen.go --model User --methods Get,Pagination,Create

type User struct {
	ID        string    `gorm:"default:generated();"`
	CreatedAt time.Time `gorm:"default:now();"`
	DeletedAt gorm.DeletedAt
}

func (u *User) into() models.User {
	return models.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		DeletedAt: sql.NullTime(u.DeletedAt),
	}
}

func (u *User) copy(user *models.User) {
	if user == nil {
		return
	}
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.DeletedAt = gorm.DeletedAt(user.DeletedAt)
}

func (d *userDB) LastApprovedUserForm(user models.User) (models.UserForm, error) {
	form := UserForm{}
	err := d.db.
		Order("created_at desc").
		Take(&form, "user_id = ? AND state = ?", user.ID, models.UserFormStateApproved).
		Error

	if err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return form.into(), nil
}

func (d *userDB) MostRelevantUserForm(user models.User) (models.UserForm, error) {
	form := UserForm{}
	err := approvedOrFirstUserFormFilter(d.db).Take(&form, "user_id = ?", user.ID).Error
	if err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return form.into(), nil
}

func (d *userDB) filter(query *gorm.DB, config *models.UsersFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.ID) > 0 {
		query = query.Where("id IN ?", config.ID)
	}

	return query
}
