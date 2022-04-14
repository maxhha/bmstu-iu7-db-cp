package ports

import "errors"

var ErrNoRole = errors.New("no role")
var ErrRecordNotFound = errors.New("not found")

var ErrUserFormIsNil = errors.New("user form is nil")
var ErrProductIsNil = errors.New("product is nil")
var ErrAuctionIsNil = errors.New("auction is nil")

var ErrInvalidFirst = errors.New("first must be positive")
