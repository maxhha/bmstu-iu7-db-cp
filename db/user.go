package db

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"default:generated();" json:"id"`
	CreatedAt    time.Time `gorm:"default:now();"`
	DeletedAt    sql.NullTime
	BlockedUntil sql.NullTime
}

func (u *User) LastApprovedUserForm(db *gorm.DB) (UserForm, error) {
	form := UserForm{}
	err := db.
		Order("created_at desc").
		Take(&form, "user_id = ? AND state = 'APPROVED'", u.ID).
		Error

	return form, err
}

func (u *User) MostRelevantUserForm(db *gorm.DB) (UserForm, error) {
	form := UserForm{}
	err := form.MostRelevantFilter(db).
		Take(&form, "user_id = ?", u.ID).
		Error

	return form, err
}
