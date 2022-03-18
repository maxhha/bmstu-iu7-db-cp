package graph

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/model"
	"auction-back/jwt"
	"auction-back/test"
	"context"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func init() {
	os.Setenv("SIGNING_KEY", "test")
	jwt.Init()
}

type RegisterSuite struct {
	test.GraphSuite
	resolver *Resolver
}

func (s *RegisterSuite) SetupTest() {
	s.GraphSuite.SetupTest()
	s.resolver = New(s.DB, &s.TokenMock)
}

func (s *RegisterSuite) TearDownTest() {
	s.GraphSuite.TearDownTest()
}

func (s *RegisterSuite) TestRegister() {
	id := "user-test"

	s.SqlMock.ExpectQuery("INSERT INTO \"users\"").
		WithArgs(nil, nil).
		WillReturnRows(test.MockRows(db.User{ID: id, CreatedAt: time.Now()}))

	result, err := s.resolver.Mutation().Register(context.Background())
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)

	token_id, err := jwt.ParseUser(result.Token)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), token_id)

	require.Equal(s.T(), id, *token_id)
}

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}

type ApproveUserEmailSuite struct {
	test.GraphSuite
	resolver *Resolver
}

func (s *ApproveUserEmailSuite) SetupTest() {
	s.GraphSuite.SetupTest()
	s.resolver = New(s.DB, &s.TokenMock)
}

func (s *ApproveUserEmailSuite) TearDownTest() {
	s.GraphSuite.TearDownTest()
}

func (s *ApproveUserEmailSuite) TestApproveUserEmail() {
	token := "123456"
	email := "email-test"
	viewer := db.User{ID: "user-test"}
	user_form := db.UserForm{ID: "test"}

	s.TokenMock.
		On("Activate", db.TokenActionApproveUserEmail, token, &viewer).
		Return(
			db.Token{Data: map[string]interface{}{"email": email}},
			nil,
		)

	s.SqlMock.ExpectQuery("SELECT \\* FROM \"user_forms\"").
		WithArgs(viewer.ID).
		WillReturnError(gorm.ErrRecordNotFound)

	s.SqlMock.ExpectQuery("INSERT INTO \"user_forms\"").
		WillReturnRows(test.MockRows(user_form))

	s.SqlMock.ExpectExec("UPDATE \"user_forms\" SET \"email\"").
		WithArgs(email, sqlmock.AnyArg(), user_form.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := auth.WithViewer(context.Background(), &viewer)
	result, err := s.resolver.Mutation().ApproveUserEmail(ctx, &model.TokenInput{Token: token})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result.User, &viewer)
}

func TestApproveUserEmailSuite(t *testing.T) {
	suite.Run(t, new(ApproveUserEmailSuite))
}
