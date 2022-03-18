package test

import (
	"auction-back/db"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TokenService struct {
	mock.Mock
}

func (m *TokenService) Validate(action db.TokenAction, data map[string]interface{}) error {
	args := m.Called(action, data)
	return args.Error(0)
}

func (m *TokenService) Send(token db.Token) error {
	args := m.Called(token)

	return args.Error(0)
}

type GraphSuite struct {
	suite.Suite
	SqlDB     *sql.DB
	DB        *gorm.DB
	SqlMock   sqlmock.Sqlmock
	TokenMock TokenService
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
}

func (s *GraphSuite) TearDownTest() {
	s.SqlDB.Close()
}
