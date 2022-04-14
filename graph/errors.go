package graph

import "errors"

var ErrViewerNotOwner = errors.New("viewer is not owner")
var ErrNotEditable = errors.New("not editable")
var ErrFloatIsNotExact = errors.New("float is not exact")
var ErrAlreadyExists = errors.New("already exists")
