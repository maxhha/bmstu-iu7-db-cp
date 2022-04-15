package ports

import (
	"errors"
	"fmt"
)

var ErrNoRole = errors.New("no role")
var ErrRecordNotFound = errors.New("not found")

var ErrIsNil = errors.New("is nil")
var ErrUserFormIsNil = fmt.Errorf("user form %w", ErrIsNil)
var ErrProductIsNil = fmt.Errorf("product %w", ErrIsNil)
var ErrAuctionIsNil = fmt.Errorf("auction %w", ErrIsNil)

var ErrInvalidFirst = errors.New("first must be positive")
