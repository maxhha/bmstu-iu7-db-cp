package graph

import (
	"auction-back/db"
	"auction-back/jwt"
	"auction-back/test"
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CreateUserSuite struct {
	suite.Suite
	db       *sql.DB
	DB       *gorm.DB
	mock     sqlmock.Sqlmock
	resolver *Resolver
}

func init() {
	os.Setenv("SIGNING_KEY", "test")
	jwt.Init()
}

func (s *CreateUserSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.db)
	require.NotNil(s.T(), s.mock)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 s.db,
		PreferSimpleProtocol: true,
	})

	s.DB, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(s.T(), err)

	s.resolver = New(s.DB)
}

func (s *CreateUserSuite) TearDownTest() {
	s.db.Close()
}

func (s *CreateUserSuite) TestCreateUser() {
	id := "user-test"

	s.mock.ExpectQuery("INSERT INTO \"users\"").
		WithArgs(nil, nil).
		WillReturnRows(test.MockRows(db.User{ID: id, CreatedAt: time.Now()}))

	result, err := s.resolver.Mutation().CreateUser(context.Background())
	require.Nil(s.T(), err)
	require.NotNil(s.T(), result)

	token_id, err := jwt.ParseUser(result.Token)
	require.Nil(s.T(), err)
	require.NotNil(s.T(), token_id)

	require.Equal(s.T(), id, *token_id)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(CreateUserSuite))
}
