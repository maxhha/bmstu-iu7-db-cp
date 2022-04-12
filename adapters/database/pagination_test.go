package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PaginationQueryByCreatedAtDescSuite struct {
	DatabaseSuite
}

func (s *PaginationQueryByCreatedAtDescSuite) TestGetAllQuery() {
	query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&Account{}).Where("user_id = ?", "USER_ID")
		query, err := paginationQueryByCreatedAtDesc(tx, nil, nil)
		require.NoError(s.T(), err)
		return query.Find(&Account{})
	})

	assert.Equal(
		s.T(),
		`SELECT * FROM "accounts" WHERE user_id = 'USER_ID' ORDER BY created_at desc`,
		query,
	)
}

func (s *PaginationQueryByCreatedAtDescSuite) TestGetFirstNQuery() {
	first := 1234

	query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&Account{}).Where("user_id = ?", "USER_ID")
		query, err := paginationQueryByCreatedAtDesc(tx, &first, nil)
		require.NoError(s.T(), err)
		return query.Find(&Account{})
	})

	assert.Equal(
		s.T(),
		s.SQL(`
			SELECT *
			FROM "accounts"
			WHERE user_id = 'USER_ID'
			ORDER BY created_at desc LIMIT %d
		`, first+1),
		query,
	)
}

func (s *PaginationQueryByCreatedAtDescSuite) TestGetAfterQuery() {
	after := "test-account"

	query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&Account{}).Where("user_id = ?", "USER_ID")
		query, err := paginationQueryByCreatedAtDesc(tx, nil, &after)
		require.NoError(s.T(), err)
		return query.Find(&Account{})
	})

	assert.Equal(
		s.T(),
		s.SQL(`
			SELECT *
			FROM "accounts"
			WHERE user_id = 'USER_ID' 
			AND created_at < ANY(
				SELECT "created_at"
				FROM "accounts"
				WHERE user_id = 'USER_ID'
				AND id = '%s'
			)
			ORDER BY created_at desc
		`, after),
		query,
	)
}

func TestPaginationQueryByCreatedAtDescSuite(t *testing.T) {
	suite.Run(t, new(PaginationQueryByCreatedAtDescSuite))
}
