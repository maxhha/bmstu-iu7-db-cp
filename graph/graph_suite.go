package graph

import (
	"auction-back/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GraphSuite struct {
	test.DBSuite
	TokenMock test.TokenPort
	BankMock  test.BankPort
	resolver  *Resolver
}

func (s *GraphSuite) SetupTest() {
	var err error
	s.SqlDB, s.SqlMock, err = sqlmock.New()
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.SqlDB)
	require.NotNil(s.T(), s.SqlMock)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 s.SqlDB,
		PreferSimpleProtocol: true,
	})

	s.DB, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(s.T(), err)

	s.resolver = New(s.DB, &s.TokenMock, &s.BankMock)
}
