package graph

import (
	"auction-back/auth"
	"auction-back/models"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestX(t *testing.T) {

}

type RequestModerateUserFormSuite struct {
	GraphSuite
}

func (s *RequestModerateUserFormSuite) TestEditableForm() {
	viewer := models.User{ID: "user-test"}
	phone := "phone"
	name := "name"
	password := "password"
	email := "email"

	for _, state := range models.AllUserFormState {
		form := models.UserForm{
			State:    state,
			Phone:    &phone,
			Name:     &name,
			Password: &password,
			Email:    &email,
		}
		if !form.IsEditable() {
			continue
		}
		ctx := auth.WithViewer(context.Background(), viewer)

		s.DB.UserFormMock.On("Take", mock.Anything).Return(form, nil).Once()
		s.TokenMock.
			On("Create", models.TokenActionModerateUserForm, viewer, map[string]interface{}{}).
			Return(nil)

		result, err := s.resolver.Mutation().RequestModerateUserForm(ctx)

		require.NoError(
			s.T(),
			err,
			"should not return error for form with state %s",
			state,
		)
		require.Equal(s.T(), result, true)
	}
}

func (s *RequestModerateUserFormSuite) TestModeratingForm() {
	viewer := models.User{ID: "user-test"}
	phone := "phone"
	name := "name"
	password := "password"
	email := "email"

	form := models.UserForm{
		State:    models.UserFormStateModerating,
		Phone:    &phone,
		Name:     &name,
		Password: &password,
		Email:    &email,
	}

	ctx := auth.WithViewer(context.Background(), viewer)
	s.DB.UserFormMock.On("Take", mock.Anything).Return(form, nil)

	result, err := s.resolver.Mutation().RequestModerateUserForm(ctx)
	require.ErrorIs(s.T(), err, ErrUserFormModerating)
	require.Equal(s.T(), result, false)

}

func TestRequestModerateUserFormSuite(t *testing.T) {
	suite.Run(t, new(RequestModerateUserFormSuite))
}

type ApproveModerateUserFormSuite struct {
	GraphSuite
}

func (s *ApproveModerateUserFormSuite) TestApproveModerateUserForm() {
	token := "123456"
	viewer := models.User{ID: "user-test"}
	user_form := models.UserForm{ID: "form-test"}

	ctx := auth.WithViewer(context.Background(), viewer)

	s.TokenMock.
		On("Activate", models.TokenActionModerateUserForm, token, viewer).
		Return(models.Token{UserID: viewer.ID}, nil)

	s.DB.UserFormMock.On("Take", mock.Anything).Return(user_form, nil)

	s.DB.UserFormMock.On("Update", mock.MatchedBy(func(form *models.UserForm) bool {
		return form.State == models.UserFormStateModerating
	})).Return(nil)

	result, err := s.resolver.Mutation().ApproveModerateUserForm(ctx, models.TokenInput{Token: token})
	require.NoError(s.T(), err)
	require.Equal(s.T(), result.User, &viewer)
}

func TestApproveModerateUserFormSuite(t *testing.T) {
	suite.Run(t, new(ApproveModerateUserFormSuite))
}
