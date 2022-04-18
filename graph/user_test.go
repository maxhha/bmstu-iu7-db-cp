package graph

import (
	"auction-back/models"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

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
