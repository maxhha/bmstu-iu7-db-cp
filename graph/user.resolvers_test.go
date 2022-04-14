package graph

import (
	"auction-back/auth"
	"auction-back/jwt"
	"auction-back/models"
	"auction-back/ports"
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

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
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

type LoginSuite struct {
	GraphSuite
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (s *LoginSuite) Test() {
	userID := "test-user"
	password := "12345"
	passwordHash, err := hashPassword(password)
	require.NoError(s.T(), err)
	passwordHash2 := "password-hash"
	ctx := context.Background()
	input := models.LoginInput{Password: password}

	cases := []struct {
		Name                 string
		GetLoginFormUserForm models.UserForm
		GetLoginFormError    error
		Error                error
	}{
		{
			Name:                 "Fail find user form",
			GetLoginFormUserForm: models.UserForm{},
			GetLoginFormError:    ports.ErrRecordNotFound,
			Error:                ports.ErrRecordNotFound,
		},
		{
			Name:                 "Form without password",
			GetLoginFormUserForm: models.UserForm{UserID: userID},
			Error:                ErrNoPassword,
		},
		{
			Name:                 "Enter wrong password",
			GetLoginFormUserForm: models.UserForm{UserID: userID, Password: &passwordHash2},
			Error:                ErrPasswordMissmatch,
		},
		{
			Name: "Success",
			GetLoginFormUserForm: models.UserForm{
				UserID:   userID,
				Password: &passwordHash,
			},
		},
	}

	for _, c := range cases {
		s.DB.UserFormMock.On("GetLoginForm", mock.Anything).
			Return(c.GetLoginFormUserForm, c.GetLoginFormError).
			Once()

		result, err := s.resolver.Mutation().Login(ctx, input)

		if c.Error == nil {
			require.NoError(s.T(), err, "error not exists in case '%s'", c.Name)
			require.NotNil(s.T(), result, "result must be in case '%s'", c.Name)

			uid, err := jwt.ParseUser(result.Token)
			require.NoError(
				s.T(),
				err,
				"must parse given token in case '%s'",
				c.Name,
			)
			require.Equal(
				s.T(),
				userID,
				uid,
				"token must match user id in case '%s'",
				c.Name,
			)
		} else {
			require.ErrorIs(
				s.T(),
				err,
				c.Error,
				"error must exists in case '%s'",
				c.Name,
			)
		}
	}
}

type ApproveSetUserEmailSuite struct {
	GraphSuite
}

func TestApproveSetUserEmailSuite(t *testing.T) {
	suite.Run(t, new(ApproveSetUserEmailSuite))
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

	s.DB.UserFormMock.On("Take", mock.Anything).Return(models.UserForm{}, ports.ErrRecordNotFound)

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

type UpdateUserPasswordSuite struct {
	GraphSuite
}

func TestUpdateUserPasswordSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserPasswordSuite))
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

type ViewerSuite struct {
	GraphSuite
}

func TestViewerSuite(t *testing.T) {
	suite.Run(t, new(ViewerSuite))
}

func (s *ViewerSuite) TestViewer() {
	viewer := models.User{ID: "user-test"}
	ctx := auth.WithViewer(context.Background(), viewer)

	result, err := s.resolver.Query().Viewer(ctx)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), result, &viewer)
}

type UserRequestSetUserEmail struct {
	GraphSuite
}

func TestUserRequestSetUserEmail(t *testing.T) {
	suite.Run(t, new(UserRequestSetUserEmail))
}

func (s *UserRequestSetUserEmail) Test() {
	viewer := models.User{ID: "test-user"}
	email := "new-email"
	cases := []struct {
		Name    string
		Context context.Context
		Mock    func()
		Error   error
	}{
		{
			Name:    "No viewer",
			Context: context.Background(),
			Mock:    func() {},
			Error:   auth.ErrUnauthorized,
		},
		{
			Name:    "Fail create token",
			Context: auth.WithViewer(context.Background(), viewer),
			Mock: func() {
				s.TokenMock.On(
					"Create",
					models.TokenActionSetUserEmail,
					viewer,
					map[string]interface{}{"email": email},
				).Return(sql.ErrConnDone).Once()
			},
			Error: sql.ErrConnDone,
		},
		{
			Name:    "Success",
			Context: auth.WithViewer(context.Background(), viewer),
			Mock: func() {
				s.TokenMock.On(
					"Create",
					models.TokenActionSetUserEmail,
					viewer,
					map[string]interface{}{"email": email},
				).Return(nil).Once()
			},
		},
	}

	for _, c := range cases {
		c.Mock()
		ok, err := s.resolver.Mutation().RequestSetUserEmail(c.Context, models.RequestSetUserEmailInput{Email: email})

		if c.Error == nil {
			require.NoError(s.T(), err, "[%s] should not have error", c.Name)
			require.Equal(s.T(), ok, true)
		} else {
			require.ErrorIs(
				s.T(),
				err,
				c.Error,
				"[%s] should have error",
				c.Name,
			)
			require.Equal(s.T(), ok, false)
		}
	}
}

type UserRequestSetUserPhone struct {
	GraphSuite
}

func TestUserRequestSetUserPhone(t *testing.T) {
	suite.Run(t, new(UserRequestSetUserPhone))
}

func (s *UserRequestSetUserPhone) Test() {
	viewer := models.User{ID: "test-user"}
	phone := "new-phone"
	cases := []struct {
		Name    string
		Context context.Context
		Mock    func()
		Error   error
	}{
		{
			Name:    "No viewer",
			Context: context.Background(),
			Mock:    func() {},
			Error:   auth.ErrUnauthorized,
		},
		{
			Name:    "Fail create token",
			Context: auth.WithViewer(context.Background(), viewer),
			Mock: func() {
				s.TokenMock.On(
					"Create",
					models.TokenActionSetUserPhone,
					viewer,
					map[string]interface{}{"phone": phone},
				).Return(sql.ErrConnDone).Once()
			},
			Error: sql.ErrConnDone,
		},
		{
			Name:    "Success",
			Context: auth.WithViewer(context.Background(), viewer),
			Mock: func() {
				s.TokenMock.On(
					"Create",
					models.TokenActionSetUserPhone,
					viewer,
					map[string]interface{}{"phone": phone},
				).Return(nil).Once()
			},
		},
	}

	for _, c := range cases {
		c.Mock()
		ok, err := s.resolver.Mutation().RequestSetUserPhone(c.Context, models.RequestSetUserPhoneInput{Phone: phone})

		if c.Error == nil {
			require.NoError(s.T(), err, "[%s] should not have error", c.Name)
			require.Equal(s.T(), ok, true)
		} else {
			require.ErrorIs(
				s.T(),
				err,
				c.Error,
				"[%s] should have error",
				c.Name,
			)
			require.Equal(s.T(), ok, false)
		}
	}
}
