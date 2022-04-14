package ports

import "errors"

var ErrNoRole = errors.New("no role")
var ErrUserFormIsNil = errors.New("user form is nil")
var ErrInvalidFirst = errors.New("first must be positive")
