package graph

import (
	"auction-back/auth"
	"auction-back/jwt"
	"auction-back/models"
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

	s.DB.UserMock.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = id
	})

	result, err := s.resolver.Mutation().Register(context.Background())
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)

	token_id, err := jwt.ParseUser(result.Token)
	require.NoError(s.T(), err)
	require.Equal(s.T(), id, token_id)
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
		On("Activate", models.TokenActionSetUserEmail, token, viewer).
		Return(
			models.Token{Data: map[string]interface{}{"email": email}},
			nil,
		)

	s.DB.UserFormMock.On("Take", mock.Anything).Return(models.UserForm{}, sql.ErrNoRows)

	s.DB.UserFormMock.On("Create", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		form := args.Get(0).(*models.UserForm)
		*form = user_form
	})

	s.DB.UserFormMock.On("Update", mock.MatchedBy(func(form *models.UserForm) bool {
		return *form.Email == email
	})).Return(nil)

	ctx := auth.WithViewer(context.Background(), viewer)
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

	ctx := auth.WithViewer(context.Background(), viewer)

	s.DB.UserFormMock.On("Take", mock.Anything).Return(user_form, nil)

	hash, err := hashPassword(password)
	require.NoError(s.T(), err)

	s.DB.UserFormMock.On("Update", mock.MatchedBy(func(form *models.UserForm) bool {
		return *form.Password == hash
	})).Return(nil)

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
	ctx := auth.WithViewer(context.Background(), viewer)

	result, err := s.resolver.Query().Viewer(ctx)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result, &viewer)
}

func TestViewerSuite(t *testing.T) {
	suite.Run(t, new(ViewerSuite))
}
