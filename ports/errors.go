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
var ErrAccountIsNil = fmt.Errorf("account %w", ErrIsNil)
var ErrOfferIsNil = fmt.Errorf("offer %w", ErrIsNil)
var ErrTransactionIsNil = fmt.Errorf("transaction %w", ErrIsNil)
var ErrBankIsNil = fmt.Errorf("bank %w", ErrIsNil)

var ErrInvalidFirst = errors.New("first must be positive")
