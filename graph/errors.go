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
var ErrAuctionIsNotStarted = errors.New("auction is not started")
var ErrAuctionIsFinished = errors.New("auction is finished")
var ErrAccountTypeNotUser = errors.New("account type is not USER")
var ErrNoCurrency = errors.New("no currency")
var ErrNotAnoughMoney = errors.New("not enough money")
