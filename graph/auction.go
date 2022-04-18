package graph

import (
	"auction-back/models"
)

func IsAuctionOwner(viewer models.User, auction models.Auction) error {
	if auction.SellerID != viewer.ID {
		return ErrViewerNotOwner
	}

	return nil
}
