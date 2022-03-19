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
	os.Setenv("PASSWORD_HASH_SALT", "test")
	jwt.Init()
	InitPasswordHashSalt()
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

type ApproveSetUserEmailSuite struct {
	test.GraphSuite
	resolver *Resolver
}

func (s *ApproveSetUserEmailSuite) SetupTest() {
	s.GraphSuite.SetupTest()
	s.resolver = New(s.DB, &s.TokenMock)
}

func (s *ApproveSetUserEmailSuite) TearDownTest() {
	s.GraphSuite.TearDownTest()
}

func (s *ApproveSetUserEmailSuite) TestApproveSetUserEmail() {
	token := "123456"
	email := "email-test"
	viewer := db.User{ID: "user-test"}
	user_form := db.UserForm{ID: "test"}

	s.TokenMock.
		On("Activate", db.TokenActionSetUserEmail, token, &viewer).
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
	result, err := s.resolver.Mutation().ApproveSetUserEmail(ctx, &model.TokenInput{Token: token})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result.User, &viewer)
}

func TestApproveSetUserEmailSuite(t *testing.T) {
	suite.Run(t, new(ApproveSetUserEmailSuite))
}

type UpdateUserPasswordSuite struct {
	test.GraphSuite
	resolver *Resolver
}

func (s *UpdateUserPasswordSuite) SetupTest() {
	s.GraphSuite.SetupTest()
	s.resolver = New(s.DB, &s.TokenMock)
}

func (s *UpdateUserPasswordSuite) TearDownTest() {
	s.GraphSuite.TearDownTest()
}

func (s *UpdateUserPasswordSuite) TestUpdatePassword() {
	password := "test-password"
	viewer := db.User{ID: "user-test"}
	user_form := db.UserForm{
		ID:    "test",
		State: db.UserFormStateCreated,
	}

	ctx := auth.WithViewer(context.Background(), &viewer)

	s.SqlMock.ExpectQuery("SELECT \\* FROM \"user_forms\"").
		WithArgs(viewer.ID).
		WillReturnRows(test.MockRows(user_form))

	hash, err := hashPassword(password)
	require.NoError(s.T(), err)

	s.SqlMock.ExpectExec("UPDATE \"user_forms\" SET \"password\"").
		WithArgs(hash, sqlmock.AnyArg(), user_form.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := s.resolver.Mutation().UpdateUserPassword(ctx, &model.UpdateUserPasswordInput{Password: password})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result.User, &viewer)
}

func TestUpdateUserPasswordSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserPasswordSuite))
}
