package model

import "auction-back/db"

type User struct {
	ID string `json:"id"`
	DB *db.User
}

type Balance struct {
	Available float64 `json:"available"`
	DB        *db.User
}

func (u *User) From(user *db.User) (*User, error) {
	u.ID = user.ID
	u.DB = user

	return u, nil
}

func (b *Balance) From(user *db.User) (*Balance, error) {
	// b.Available = user.Available
	b.DB = user

	return b, nil
}
