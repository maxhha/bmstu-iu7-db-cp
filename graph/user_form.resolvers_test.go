package graph

import (
	"auction-back/auth"
	"auction-back/models"
	"auction-back/test"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RequestModerateUserFormSuite struct {
	GraphSuite
}

func (s *RequestModerateUserFormSuite) TestRequestModerateUserForm() {
	viewer := models.User{ID: "user-test"}
	phone := "phone"
	name := "name"
	password := "password"
	email := "email"

	draft_form := models.UserForm{
		State:    models.UserFormStateCreated,
		Phone:    &phone,
		Name:     &name,
		Password: &password,
		Email:    &email,
	}

	ctx := auth.WithViewer(context.Background(), &viewer)

	s.SqlMock.ExpectQuery("SELECT \\* FROM \"user_forms\"").
		WithArgs(viewer.ID).
		WillReturnRows(test.MockRows(draft_form))

	s.TokenMock.
		On("Create", models.TokenActionModerateUserForm, &viewer, map[string]interface{}{}).
		Return(nil)

	result, err := s.resolver.Mutation().RequestModerateUserForm(ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), result, true)
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

	ctx := auth.WithViewer(context.Background(), &viewer)

	s.TokenMock.
		On("Activate", models.TokenActionModerateUserForm, token, &viewer).
		Return(models.Token{UserID: viewer.ID}, nil)

	s.SqlMock.ExpectQuery("SELECT \\* FROM \"user_forms\"").
		WithArgs(viewer.ID, models.UserFormStateCreated, models.UserFormStateDeclained).
		WillReturnRows(test.MockRows(user_form))

	s.SqlMock.ExpectExec("UPDATE \"user_forms\" SET \"state\"").
		WithArgs(models.UserFormStateModerating, sqlmock.AnyArg(), user_form.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := s.resolver.Mutation().ApproveModerateUserForm(ctx, models.TokenInput{Token: token})
	require.NoError(s.T(), err)
	require.Equal(s.T(), result.User, &viewer)
}

func TestApproveModerateUserFormSuite(t *testing.T) {
	suite.Run(t, new(ApproveModerateUserFormSuite))
}
