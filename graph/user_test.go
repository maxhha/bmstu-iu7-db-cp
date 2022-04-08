package graph

import (
	"testing"
)

func TestXxx(t *testing.T) {

}

// type GetOrCreateUserDraftFormSuite struct {
// 	test.DBSuite
// }

// func (s *GetOrCreateUserDraftFormSuite) TestClone() {
// 	viewer := models.User{ID: "test-user"}
// 	approved_form := models.UserForm{
// 		ID:     "approved-form",
// 		State:  models.UserFormStateApproved,
// 		UserID: viewer.ID,
// 		Model: gorm.Model{
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}

// 	s.SqlMock.ExpectQuery("SELECT \\* FROM \"user_forms\" WHERE user_id = ").
// 		WithArgs(viewer.ID).
// 		WillReturnRows(test.MockRows(approved_form))

// 	s.SqlMock.ExpectQuery("INSERT INTO \"user_forms\"").
// 		WithArgs(
// 			test.AfterTime{Time: approved_form.CreatedAt},
// 			test.AfterTime{Time: approved_form.UpdatedAt},
// 			nil,
// 			approved_form.UserID,
// 			models.UserFormStateCreated,
// 			nil,
// 			nil,
// 			nil,
// 			nil,
// 			nil,
// 		).
// 		WillReturnRows(test.MockRows(models.UserForm{ID: "created-form"}))

// 	_, err := getOrCreateUserDraftForm(s.DB, viewer)
// 	require.NoError(s.T(), err)
// }

// func TestGetOrCreateUserDraftFormSuite(t *testing.T) {
// 	suite.Run(t, new(GetOrCreateUserDraftFormSuite))
// }
