package graph

import (
	"auction-back/db"
	"auction-back/jwt"
	"auction-back/test"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateUserSuite struct {
	test.GraphSuite
	resolver *Resolver
}

func init() {
	os.Setenv("SIGNING_KEY", "test")
	jwt.Init()
}

func (s *CreateUserSuite) SetupTest() {
	s.GraphSuite.SetupTest()
	s.resolver = New(s.DB, &s.TokenMock)
}

func (s *CreateUserSuite) TearDownTest() {
	s.GraphSuite.TearDownTest()
}

func (s *CreateUserSuite) TestCreateUser() {
	id := "user-test"

	s.SqlMock.ExpectQuery("INSERT INTO \"users\"").
		WithArgs(nil, nil).
		WillReturnRows(test.MockRows(db.User{ID: id, CreatedAt: time.Now()}))

	result, err := s.resolver.Mutation().CreateUser(context.Background())
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)

	token_id, err := jwt.ParseUser(result.Token)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), token_id)

	require.Equal(s.T(), id, *token_id)
}

func TestCreateUserSuite(t *testing.T) {
	suite.Run(t, new(CreateUserSuite))
}
