package model

import "auction-back/db"

type User struct {
	ID      string   `json:"id"`
	Balance *Balance `json:"balance"`
	DB      *db.User
}

func (u *User) From(user *db.User) (*User, error) {
	u.ID = user.ID
	u.Balance = &Balance{
		Available: user.Available,
		Blocked:   0,
	}
	u.DB = user

	return u, nil
}
