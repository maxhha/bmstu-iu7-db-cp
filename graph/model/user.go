package model

import "auction-back/db"

func (u *User) From(user *db.User) (*User, error) {
	u.ID = user.ID
	u.Balance = &Balance{
		Available: user.Available,
		Blocked:   0,
	}

	return u, nil
}
