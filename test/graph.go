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

func (t *TokenService) Create(action db.TokenAction, viewer *db.User, data map[string]interface{}) error {
	args := t.Called(action, viewer, data)
	return args.Error(0)
}

func (t *TokenService) Activate(action db.TokenAction, token_code string, viewer *db.User) (db.Token, error) {
	args := t.Called(action, token_code, viewer)
	return args.Get(0).(db.Token), args.Error(1)
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
