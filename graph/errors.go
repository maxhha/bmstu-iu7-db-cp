package graph

import (
	"auction-back/ports"
	"errors"
	"fmt"
)

var ErrViewerNotOwner = errors.New("viewer is not owner")
var ErrNotEditable = errors.New("not editable")
var ErrAlreadyExists = errors.New("already exists")
var ErrCurrencyIsNil = fmt.Errorf("currency %w", ports.ErrIsNil)
