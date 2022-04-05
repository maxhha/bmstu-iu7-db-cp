package graph

import (
	"auction-back/auth"
	"auction-back/jwt"
	"auction-back/models"
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
	GraphSuite
}

func (s *RegisterSuite) TestRegister() {
	id := "user-test"

	s.SqlMock.ExpectQuery("INSERT INTO \"users\"").
		WithArgs(nil, nil).
		WillReturnRows(test.MockRows(models.User{ID: id, CreatedAt: time.Now()}))

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
	GraphSuite
}

func (s *ApproveSetUserEmailSuite) TestApproveSetUserEmail() {
	token := "123456"
	email := "email-test"
	viewer := models.User{ID: "user-test"}
	user_form := models.UserForm{ID: "test"}

	s.TokenMock.
		On("Activate", models.TokenActionSetUserEmail, token, &viewer).
		Return(
			models.Token{Data: map[string]interface{}{"email": email}},
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
	result, err := s.resolver.Mutation().ApproveSetUserEmail(ctx, models.TokenInput{Token: token})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result.User, &viewer)
}

func TestApproveSetUserEmailSuite(t *testing.T) {
	suite.Run(t, new(ApproveSetUserEmailSuite))
}

type UpdateUserPasswordSuite struct {
	GraphSuite
}

func (s *UpdateUserPasswordSuite) TestUpdatePassword() {
	password := "test-password"
	viewer := models.User{ID: "user-test"}
	user_form := models.UserForm{
		ID:    "test",
		State: models.UserFormStateCreated,
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

	result, err := s.resolver.Mutation().UpdateUserPassword(ctx, models.UpdateUserPasswordInput{Password: password})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result.User, &viewer)
}

func TestUpdateUserPasswordSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserPasswordSuite))
}

type ViewerSuite struct {
	GraphSuite
}

func (s *ViewerSuite) TestViewer() {
	viewer := models.User{ID: "user-test"}
	ctx := auth.WithViewer(context.Background(), &viewer)

	result, err := s.resolver.Query().Viewer(ctx)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result, &viewer)
}

func TestViewerSuite(t *testing.T) {
	suite.Run(t, new(ViewerSuite))
}
