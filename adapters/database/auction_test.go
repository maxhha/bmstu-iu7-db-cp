package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuctionsFinishAuctionsQuerySuite struct {
	DatabaseSuite
}

func TestAuctionsFinishAuctionsQuerySuite(t *testing.T) {
	suite.Run(t, new(AuctionsFinishAuctionsQuerySuite))
}

func (s *AuctionsFinishAuctionsQuerySuite) TestSimpleQuery() {
	query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var objs []Auction
		adb := s.database.Auction().(*auctionDB)
		query := tx.Model(&Auction{})
		query = adb.forFinishAuctionsQuery(query, "3 mins", "10 mins")
		return query.Find(&objs)
	})

	assert.Equal(
		s.T(),
		s.SQL(`
			SELECT * FROM "auctions" 
			WHERE state = 'STARTED' 
			AND NOW() - interval '3 mins' > (
				SELECT COALESCE(MAX(offers.created_at), auctions.started_at)
				FROM offers
				WHERE offers.auction_id = auctions.id
			)
			AND NOW() > COALESCE(auction.scheduled_finish_at, auctions.started_at + interval '10 mins')
			FOR UPDATE OF "auctions"
		`),
		query,
	)
}
