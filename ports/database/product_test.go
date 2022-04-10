package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ProductOwnersQuerySuite struct {
	DatabaseSuite
}

func TestProductOwnersQuerySuite(t *testing.T) {
	suite.Run(t, new(ProductOwnersQuerySuite))
}

func (s *ProductOwnersQuerySuite) TestSimpleQuery() {
	query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		query := s.DB.Model(&Product{})
		return s.database.Product().(*productDB).ownersQuery(query)
	})

	assert.Equal(
		s.T(),
		s.SQL(`
			SELECT id as product_id, creator_id as owner_id, created_at as from_date
			FROM "products"
			UNION ALL
			SELECT products.id as product_id, auctions.buyer_id as owner_id, auctions.finished_at as from_date
			FROM "products"
			FULL JOIN auctions ON auctions.product_id = products.id
		`),
		query,
	)
}

func (s *ProductOwnersQuerySuite) TestSubQuery() {
	query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		query := s.DB.Model(&Product{}).Where("id = ?", "test-product")
		return s.database.Product().(*productDB).ownersQuery(query)
	})

	assert.Equal(
		s.T(),
		s.SQL(`
			SELECT id as product_id, creator_id as owner_id, created_at as from_date
			FROM "products"
			WHERE id = 'test-product'
			UNION ALL
			SELECT products.id as product_id, auctions.buyer_id as owner_id, auctions.finished_at as from_date
			FROM "products"
			FULL JOIN auctions ON auctions.product_id = products.id 
			WHERE id = 'test-product'
		`),
		query,
	)
}
