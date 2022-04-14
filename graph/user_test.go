package graph

import (
	"auction-back/models"
	"auction-back/ports"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestXxx(t *testing.T) {

}

type UserIsOwnerOrManagerSuite struct {
	GraphSuite
}

func TestUserIsOwnerOrManagerSuite(t *testing.T) {
	suite.Run(t, new(UserIsOwnerOrManagerSuite))
}

func (s *UserIsOwnerOrManagerSuite) Test() {
	cases := []struct {
		Name   string
		User   *models.User
		Viewer models.User
		Errors []error
		Mock   func()
	}{
		{
			Name:   "User is nul",
			Viewer: models.User{ID: "test-user"},
			Mock:   func() {},
			Errors: []error{ErrUserIsNil},
		},
		{
			Name:   "Is owner",
			User:   &models.User{ID: "test-user"},
			Viewer: models.User{ID: "test-user"},
			Mock:   func() {},
		},
		{
			Name:   "Is manager",
			User:   &models.User{ID: "test-user"},
			Viewer: models.User{ID: "test-manager"},
			Mock: func() {
				s.RoleMock.On(
					"HasRole",
					models.RoleTypeManager,
					models.User{ID: "test-manager"},
				).Return(nil).Once()
			},
		},
		{
			Name:   "Not owner and manager",
			User:   &models.User{ID: "test-user"},
			Viewer: models.User{ID: "test-manager"},
			Mock: func() {
				s.RoleMock.On(
					"HasRole",
					models.RoleTypeManager,
					models.User{ID: "test-manager"},
				).Return(ports.ErrNoRole).Once()
			},
			Errors: []error{
				ports.ErrNoRole,
				ErrUserNotOwner,
			},
		},
	}

	for _, c := range cases {
		c.Mock()
		err := (&userResolver{s.resolver}).isOwnerOrManager(c.Viewer, c.User)

		if len(c.Errors) == 0 {
			require.NoError(s.T(), err, "[%s] should not have any errors", c.Name)
			continue
		}

		for _, e := range c.Errors {
			require.ErrorIs(s.T(), err, e, "[%s] should have error", c.Name)
		}
	}
}

type GetOrCreateUserDraftFormSuite struct {
	GraphSuite
}

func (s *GetOrCreateUserDraftFormSuite) TestClone() {
	viewer := models.User{ID: "test-user"}
	approved_form := models.UserForm{
		ID:        "approved-form",
		State:     models.UserFormStateApproved,
		UserID:    viewer.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.DB.UserFormMock.On("Take", mock.Anything).Return(approved_form, nil)

	s.DB.UserFormMock.On("Create", mock.MatchedBy(func(form *models.UserForm) bool {
		return form.UserID == approved_form.UserID &&
			form.CreatedAt.Equal(time.Time{}) &&
			form.UpdatedAt.Equal(time.Time{}) &&
			form.State == models.UserFormStateCreated
	})).Return(nil)

	_, err := getOrCreateUserDraftForm(&s.DB, viewer)
	require.NoError(s.T(), err)
}

func TestGetOrCreateUserDraftFormSuite(t *testing.T) {
	suite.Run(t, new(GetOrCreateUserDraftFormSuite))
}
