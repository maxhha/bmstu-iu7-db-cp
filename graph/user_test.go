package graph

import (
	"auction-back/db"
	"auction-back/test"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type GetOrCreateUserDraftFormSuite struct {
	test.DBSuite
}

func (s *GetOrCreateUserDraftFormSuite) TestClone() {
	viewer := db.User{ID: "test-user"}
	approved_form := db.UserForm{
		ID:     "approved-form",
		State:  db.UserFormStateApproved,
		UserID: viewer.ID,
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	s.SqlMock.ExpectQuery("SELECT \\* FROM \"user_forms\" WHERE user_id = ").
		WithArgs(viewer.ID).
		WillReturnRows(test.MockRows(approved_form))

	s.SqlMock.ExpectQuery("INSERT INTO \"user_forms\"").
		WithArgs(
			test.AfterTime{Time: approved_form.CreatedAt},
			test.AfterTime{Time: approved_form.UpdatedAt},
			nil,
			approved_form.UserID,
			db.UserFormStateCreated,
			nil,
			nil,
			nil,
			nil,
			nil,
		).
		WillReturnRows(test.MockRows(db.UserForm{ID: "created-form"}))

	_, err := getOrCreateUserDraftForm(s.DB, &viewer)
	require.NoError(s.T(), err)
}

func TestGetOrCreateUserDraftFormSuite(t *testing.T) {
	suite.Run(t, new(GetOrCreateUserDraftFormSuite))
}
