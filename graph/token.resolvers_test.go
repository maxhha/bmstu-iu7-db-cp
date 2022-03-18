package graph

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/model"
	"auction-back/test"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateTokenSuite struct {
	test.GraphSuite
	resolver *Resolver
}

func (s *CreateTokenSuite) SetupTest() {
	s.GraphSuite.SetupTest()
	s.resolver = New(s.DB, &s.TokenMock)
}

func (s *CreateTokenSuite) TearDownTest() {
	s.GraphSuite.TearDownTest()
}

func (s *CreateTokenSuite) TestCreateToken() {
	id := "user-test"

	s.TokenMock.On("Validate", mock.Anything, mock.Anything).Return(nil)
	s.TokenMock.On("Send", mock.Anything).Return(nil)

	s.SqlMock.ExpectQuery("INSERT INTO \"tokens\"").
		WillReturnRows(test.MockRows(db.Token{ID: 123456, CreatedAt: time.Now()}))

	ctx := auth.WithViewer(context.Background(), &db.User{ID: id})
	input := model.CreateTokenInput{
		Action: model.TokenActionEnumApproveUserEmail,
		Data:   map[string]interface{}{},
	}

	result, err := s.resolver.Mutation().CreateToken(ctx, &input)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), *result, true)
}

func (s *CreateTokenSuite) TestUnauthorized() {
	ctx := context.Background()
	input := model.CreateTokenInput{
		Action: model.TokenActionEnumApproveUserEmail,
		Data:   map[string]interface{}{},
	}

	result, err := s.resolver.Mutation().CreateToken(ctx, &input)
	require.ErrorContains(s.T(), err, "unauthorized")
	require.Nil(s.T(), result)
}

func TestCreateTokenSuite(t *testing.T) {
	suite.Run(t, new(CreateTokenSuite))
}
