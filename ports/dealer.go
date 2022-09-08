package ports

type Dealer interface {
	OwnerAccept(offerId string, userId *string, comment *string) error
	BuyerAccept(offerId string, userId *string) error
}
