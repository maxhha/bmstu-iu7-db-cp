package database

import (
	"auction-back/models"
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
			JOIN auctions ON auctions.product_id = products.id
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
			JOIN auctions ON auctions.product_id = products.id 
			WHERE id = 'test-product'
		`),
		query,
	)
}

type ProductFilterSuite struct {
	DatabaseSuite
}

func TestProductFilterSuite(t *testing.T) {
	suite.Run(t, new(ProductFilterSuite))
}

func (s *ProductFilterSuite) TestFilterOwnerIDs() {
	config := models.ProductsFilter{
		OwnerIDs: []string{"test-user"},
	}
	query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return s.database.Product().(*productDB).filter(tx.Model(&Product{}), config).Find(&Product{})
	})

	assert.Equal(
		s.T(),
		s.SQL(`
			SELECT "products"."id","products"."state","products"."title","products"."description","products"."creator_id","products"."declain_reason","products"."created_at","products"."updated_at","products"."deleted_at"
			FROM "products"
			JOIN (
				SELECT *, ROW_NUMBER() OVER(PARTITION BY ofd.product_id ORDER BY ofd.from_date DESC) as owner_n
				FROM (
					SELECT id as product_id, creator_id as owner_id, created_at as from_date
					FROM "products"
					UNION ALL
					SELECT products.id as product_id, auctions.buyer_id as owner_id, auctions.finished_at as from_date
					FROM "products"
					JOIN auctions 
					ON auctions.product_id = products.id
				) as ofd
			) ofd 
			ON products.id = ofd.product_id AND ofd.owner_n = 1 AND ofd.owner_id IN ('test-user')`),
		query,
	)
}
