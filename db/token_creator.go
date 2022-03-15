package db

type TokenCreator struct {
	ID      string
	UserId  *string
	GuestId *string
	User    *User
	Guest   *Guest
}
