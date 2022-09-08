package main

import (
	"auction-back/models"
	"auction-back/ports"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
)

type Dealer struct {
	db       ports.DB
	bank     ports.Bank
	timer    *time.Timer
	stopChan chan struct{}
	stopped  bool
	timeout  time.Duration
}

var ErrWrongState = errors.New("wrong state")

func NewDealer(db ports.DB, bank ports.Bank) *Dealer {
	ch := make(chan struct{}, 1)
	return &Dealer{
		db,
		bank,
		time.NewTimer(time.Microsecond),
		ch,
		false,
		time.Microsecond,
	}
}

func (d *Dealer) Run() {
	for {
		select {
		case <-d.timer.C:
			d.timeout = time.Duration(15) * time.Second
			if _, err := d.Process(ports.FindAndSetFinishConfig{
				DefaultDuration:      "10 mins",
				TimeGapFromLastOffer: "3 mins",
			}); err != nil {
				fmt.Printf("process: %v\n", err)
			}
			d.timer.Reset(d.timeout)
		case <-d.stopChan:
			return
		}
	}
}

func (d *Dealer) Process(config ports.FindAndSetFinishConfig) ([]models.Auction, error) {
	auctions, err := d.db.Auction().FindAndSetFinish(config)
	if err != nil {
		return nil, fmt.Errorf("d.db.Auction().FindAndSetFinish(): %w", err)
	}

	fmt.Printf("finish %d auctions\n", len(auctions))

	var errs error

	for _, auction := range auctions {
		err := d.TransferringMoney(auction)
		if err != nil {
			errs = multierror.Append(
				errs,
				fmt.Errorf(
					"d.StartDeal(auctionId=%s): %w",
					auction.ID,
					err,
				),
			)
		}
	}

	return auctions, errs
}

func (d *Dealer) TransferringMoney(auction models.Auction) error {
	topOffer, err := d.db.Auction().GetTopOffer(auction)
	if err != nil {
		if errors.Is(err, ports.ErrNoRole) {
			return nil
		}
		return fmt.Errorf("d.db.Auction().GetTopOffer: %w", err)
	}

	topOffer.State = models.OfferStateAccepted
	if err = d.db.Offer().Update(&topOffer); err != nil {
		return fmt.Errorf("d.db.Offer().Update(&topOffer): %w", err)
	}

	offers, err := d.db.Offer().Find(ports.OfferFindConfig{
		Filter: &models.OffersFilter{
			AuctionIDs: []string{topOffer.AuctionID},
			States:     []models.OfferState{models.OfferStateCreated},
		},
	})
	if err != nil {
		return fmt.Errorf("d.db.Offer().Find: %w", err)
	}

	var errs error

	for _, offer := range offers {
		offer.State = models.OfferStateCancelled
		if err := d.db.Offer().Update(&offer); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("d.db.Offer().Update(id=%s): %w", offer.ID, err))
		}
	}

	if errs != nil {
		return errs
	}

	deal := models.DealState{
		State:   models.DealStateEnumTransferringMoney,
		OfferID: topOffer.ID,
	}

	err = d.db.DealState().Create(&deal)
	if err != nil {
		return fmt.Errorf("d.db.DealState().Create: %w", err)
	}

	trs, err := d.db.Transaction().Find(ports.TransactionFindConfig{
		Filter: &models.TransactionsFilter{
			OfferIDs: []string{topOffer.ID},
		},
	})
	if err != nil {
		origErr := fmt.Errorf("d.db.Transaction().Find: %w", err)
		comment := fmt.Sprintf("TransferringMoney: %s", origErr.Error())
		err = d.TransferMoneyFailed(topOffer.ID, &comment)
		if err != nil {
			return fmt.Errorf("d.TransferMoneyFailed: %w", err)
		}
		return origErr
	}

	err = d.bank.ProcessTransactions(trs)
	if err != nil {
		origErr := fmt.Errorf("d.bank.ProcessTransactions: %w", err)
		comment := fmt.Sprintf("TransferringMoney: %s", origErr.Error())
		err = d.TransferMoneyFailed(topOffer.ID, &comment)
		if err != nil {
			return fmt.Errorf("d.TransferMoneyFailed: %w", err)
		}
		return origErr
	}

	trs, err = d.db.Transaction().Find(ports.TransactionFindConfig{
		Filter: &models.TransactionsFilter{
			AuctionIDs: []string{topOffer.AuctionID},
			States:     []models.TransactionState{models.TransactionStateCreated},
		},
	})

	if err != nil {
		return fmt.Errorf("d.db.Transaction().Find(state=created): %w", err)
	}

	for _, tr := range trs {
		tr.State = models.TransactionStateCancelled
		if err := d.db.Transaction().Update(&tr); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("d.db.Transaction().Update(state=cancel): %w", err))
		}
	}

	return errs
}

func (d *Dealer) TransferMoneyFailed(offerId string, comment *string) error {
	lastDeal, err := d.db.DealState().GetLast(offerId)
	if err != nil {
		return fmt.Errorf("d.db.DealState().GetLast: %w", err)
	}

	if lastDeal.State != models.DealStateEnumTransferringMoney {
		return fmt.Errorf("%w: must be %s, but is %s", ErrWrongState, models.DealStateEnumTransferringMoney, lastDeal.State)
	}

	deal := models.DealState{
		State:   models.DealStateEnumTransferMoneyFailed,
		OfferID: offerId,
		Comment: comment,
	}

	err = d.db.DealState().Create(&deal)
	if err != nil {
		return fmt.Errorf("d.db.DealState().Create: %w", err)
	}

	return nil
}

func (d *Dealer) TransferringProduct(offerId string, approverId *string, comment *string) error {
	lastDeal, err := d.db.DealState().GetLast(offerId)
	if err != nil {
		return fmt.Errorf("d.db.DealState().GetLast: %w", err)
	}

	if lastDeal.State != models.DealStateEnumTransferringMoney && lastDeal.State != models.DealStateEnumTransferMoneyFailed {
		return fmt.Errorf("%w: must be %s or %s, but is %s", ErrWrongState, models.DealStateEnumTransferringMoney, models.DealStateEnumTransferMoneyFailed, lastDeal.State)
	}

	deal := models.DealState{
		State:     models.DealStateEnumTransferringProduct,
		OfferID:   offerId,
		Comment:   comment,
		CreatorID: approverId,
	}

	err = d.db.DealState().Create(&deal)
	if err != nil {
		return fmt.Errorf("d.db.DealState().Create: %w", err)
	}

	return nil
}

func (d *Dealer) TransferProductFailed(offerId string, approverId *string, comment *string) error {
	lastDeal, err := d.db.DealState().GetLast(offerId)
	if err != nil {
		return fmt.Errorf("d.db.DealState().GetLast: %w", err)
	}

	if lastDeal.State != models.DealStateEnumTransferringProduct {
		return fmt.Errorf("%w: must be %s, but is %s", ErrWrongState, models.DealStateEnumTransferringProduct, lastDeal.State)
	}

	deal := models.DealState{
		State:     models.DealStateEnumTransferProductFailed,
		CreatorID: approverId,
		OfferID:   offerId,
		Comment:   comment,
	}

	err = d.db.DealState().Create(&deal)
	if err != nil {
		return fmt.Errorf("d.db.DealState().Create: %w", err)
	}

	return nil
}

func (d *Dealer) Succeeded(offerId string, approverId *string, comment *string) error {
	lastDeal, err := d.db.DealState().GetLast(offerId)
	if err != nil {
		return fmt.Errorf("d.db.DealState().GetLast: %w", err)
	}

	if lastDeal.State != models.DealStateEnumTransferringProduct &&
		lastDeal.State != models.DealStateEnumTransferProductFailed &&
		lastDeal.State != models.DealStateEnumReturnMoneyFailed {
		return fmt.Errorf(
			"%w: must be %s, %s or %s, but is %s",
			ErrWrongState,
			models.DealStateEnumTransferringProduct,
			models.DealStateEnumTransferProductFailed,
			models.DealStateEnumReturnMoneyFailed,
			lastDeal.State,
		)
	}

	deal := models.DealState{
		State:     models.DealStateEnumSucceeded,
		OfferID:   offerId,
		Comment:   comment,
		CreatorID: approverId,
	}

	err = d.db.DealState().Create(&deal)
	if err != nil {
		return fmt.Errorf("d.db.DealState().Create: %w", err)
	}

	topOffer, err := d.db.Offer().Get(offerId)
	if err != nil {
		return fmt.Errorf("d.db.Offer().Get: %w", err)
	}

	auction, err := d.db.Auction().Get(topOffer.AuctionID)
	if err != nil {
		return fmt.Errorf("d.db.Auction().Get: %w", err)
	}

	topOffer.State = models.OfferStateSucceeded
	auction.BuyerID = &topOffer.UserID
	auction.State = models.AuctionStateSucceeded
	if err = d.db.Auction().Update(&auction); err != nil {
		return fmt.Errorf("d.db.Auction().Update: %w", err)
	}

	if err = d.db.Offer().Update(&topOffer); err != nil {
		return fmt.Errorf("d.db.Offer().Update: %w", err)
	}

	return nil
}

func (d *Dealer) Stop() {
	if d.stopped {
		return
	}
	d.stopped = true
	d.stopChan <- struct{}{}
}
