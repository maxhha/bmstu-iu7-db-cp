package database

import (
	"auction-back/models"
	"auction-back/ports"
	"auction-back/test"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserFormSuite struct {
	DatabaseSuite
}

func TestUserFormSuite(t *testing.T) {
	suite.Run(t, new(UserFormSuite))
}

func (s *UserFormSuite) TestGet() {
	id := "test-user-form"
	form := UserForm{
		ID: id,
	}

	cases := []struct {
		Name  string
		Mock  func()
		Error error
	}{
		{
			"Success",
			func() {
				s.SqlMock.
					ExpectQuery(s.SQL(`
						SELECT \* FROM "user_forms"
						WHERE id = \$1
						AND "user_forms"\."deleted_at" IS NULL LIMIT 1
					`)).
					WithArgs(id).
					WillReturnRows(MockRows(form))
			},
			nil,
		},
		{
			"Not found",
			func() {
				s.SqlMock.
					ExpectQuery(s.SQL(`
						SELECT \* FROM "user_forms"
						WHERE id = \$1
						AND "user_forms"\."deleted_at" IS NULL
						LIMIT 1
					`)).
					WithArgs(id).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			ports.ErrRecordNotFound,
		},
	}

	for _, c := range cases {
		c.Mock()

		userForm, err := s.database.UserForm().Get(id)
		if c.Error != nil {
			require.ErrorIs(s.T(), err, c.Error, "[%s]", c.Name)
			continue
		}

		require.NoError(s.T(), err, "[%s]", c.Name)
		require.Equal(s.T(), userForm, form.into())

	}
}

func (s *UserFormSuite) TestCreate() {
	id := "test-user-form"

	cases := []struct {
		Name  string
		Mock  func()
		Form  *models.UserForm
		Error error
	}{
		{
			"Success",
			func() {
				s.SqlMock.
					ExpectQuery(`INSERT INTO "user_forms" .* RETURNING "id"`).
					WithArgs(
						"",
						models.UserFormStateCreated,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						test.AfterTime{Time: time.Now()},
						test.AfterTime{Time: time.Now()},
						nil,
					).
					WillReturnRows(MockRows(UserForm{ID: id}))
			},
			&models.UserForm{},
			nil,
		},
		{
			"Error on create",
			func() {
				s.SqlMock.
					ExpectQuery(`INSERT INTO "user_forms" .* RETURNING "id"`).
					WithArgs(
						"",
						models.UserFormStateCreated,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						test.AfterTime{Time: time.Now()},
						test.AfterTime{Time: time.Now()},
						nil,
					).
					WillReturnError(sql.ErrConnDone)
			},
			&models.UserForm{},
			sql.ErrConnDone,
		},
		{
			"Form is nil",
			func() {},
			nil,
			ports.ErrUserFormIsNil,
		},
	}

	for _, c := range cases {
		c.Mock()

		err := s.database.UserForm().Create(c.Form)
		if c.Error != nil {
			require.ErrorIs(s.T(), err, c.Error, "[%s]", c.Name)
			continue
		}

		require.NoError(s.T(), err, "[%s]", c.Name)
		require.Equal(s.T(), c.Form.ID, id)
	}
}

func (s *UserFormSuite) TestUpdate() {
	form := UserForm{
		ID:        "test-user-form",
		UserID:    "test-user",
		State:     models.UserFormStateModerating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	form1 := form.into()
	form2 := form.into()

	cases := []struct {
		Name  string
		Mock  func()
		Form  *models.UserForm
		Error error
	}{
		{
			"Success",
			func() {
				s.SqlMock.
					ExpectExec(s.SQL(`
						UPDATE "user_forms" SET .*
						WHERE "user_forms"\."deleted_at" IS NULL AND "id" = 
					`)).
					WithArgs(
						form.UserID,
						models.UserFormStateModerating,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						form.CreatedAt,
						test.AfterTime{Time: form.UpdatedAt},
						nil,
						form.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			&form1,
			nil,
		},
		{
			"Error on update",
			func() {
				s.SqlMock.
					ExpectExec(s.SQL(`
						UPDATE "user_forms" SET .*
						WHERE "user_forms"\."deleted_at" IS NULL AND "id" =
					`)).
					WithArgs(
						form.UserID,
						models.UserFormStateModerating,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						form.CreatedAt,
						test.AfterTime{Time: form.UpdatedAt},
						nil,
						form.ID,
					).
					WillReturnError(sql.ErrConnDone)
			},
			&form2,
			sql.ErrConnDone,
		},
		{
			"Form is nil",
			func() {},
			nil,
			ports.ErrUserFormIsNil,
		},
	}

	for _, c := range cases {
		c.Mock()

		err := s.database.UserForm().Update(c.Form)
		if c.Error != nil {
			require.ErrorIs(s.T(), err, c.Error, "[%s]", c.Name)
			continue
		}

		require.NoError(s.T(), err, "[%s]", c.Name)
		require.Equal(s.T(), form.CreatedAt, c.Form.CreatedAt)
		if !c.Form.UpdatedAt.After(form.UpdatedAt) {
			s.T().Error("updated at not updated")
		}
	}
}

func (s *UserFormSuite) TestFilter() {
	baseNoOp := func(query *gorm.DB) *gorm.DB { return query.Model(&UserForm{}) }
	ids := []string{"test-id1", "test-id2"}
	idsString := "'test-id1','test-id2'"

	cases := []struct {
		Name   string
		Base   func(query *gorm.DB) *gorm.DB
		Config models.UserFormsFilter
		SQL    string
	}{
		{
			"Ids filter",
			baseNoOp,
			models.UserFormsFilter{ID: ids},
			s.SQL(`
				SELECT * FROM "user_forms"
				WHERE id IN (%s)
				AND "user_forms"."deleted_at" IS NULL
			`, idsString),
		},
		{
			"User ids filter",
			baseNoOp,
			models.UserFormsFilter{UserID: ids},
			s.SQL(`
				SELECT * FROM "user_forms"
				WHERE user_id IN (%s)
				AND "user_forms"."deleted_at" IS NULL
			`, idsString),
		},
		{
			"States filter",
			baseNoOp,
			models.UserFormsFilter{State: []models.UserFormState{
				models.UserFormStateApproved,
				models.UserFormStateDeclained,
			}},
			s.SQL(
				`
				SELECT * FROM "user_forms"
				WHERE state IN ('%s','%s')
				AND "user_forms"."deleted_at" IS NULL
				`,
				models.UserFormStateApproved,
				models.UserFormStateDeclained,
			),
		},
	}

	for _, c := range cases {
		query := s.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return (&userFormDB{s.database}).
				filter(c.Base(tx), &c.Config).
				Find(&UserForm{})
		})

		assert.Equal(s.T(), c.SQL, query, "[%s]", c.Name)
	}
}

func (s *UserFormSuite) TestPagination() {
	formsN := 10
	first := 5
	negfirst := -10
	forms := make([]UserForm, 0, formsN)

	for i := 0; i < formsN; i++ {
		forms = append(forms, UserForm{
			ID: fmt.Sprintf("test-user-form-%d", i),
		})
	}

	cases := []struct {
		Name  string
		Mock  func()
		First *int
		After *string
		Error error
		*models.PageInfo
		Forms []UserForm
	}{
		{
			"Get all",
			func() {
				s.SqlMock.
					ExpectQuery(`SELECT .* FROM "user_forms" WHERE "user_forms"\."deleted_at" IS NULL`).
					WillReturnRows(MockRows(forms))
			},
			nil,
			nil,
			nil,
			&models.PageInfo{
				StartCursor: &forms[0].ID,
				EndCursor:   &forms[len(forms)-1].ID,
			},
			forms,
		},
		{
			"Has next page",
			func() {
				s.SqlMock.
					ExpectQuery(fmt.Sprintf(
						`SELECT .* FROM "user_forms" WHERE .* LIMIT %d`,
						first+1,
					)).
					WillReturnRows(MockRows(forms[:first+1]))
			},
			&first,
			nil,
			nil,
			&models.PageInfo{
				HasNextPage: true,
				StartCursor: &forms[0].ID,
				EndCursor:   &forms[first-1].ID,
			},
			forms[:first],
		},
		{
			"Empty connection",
			func() {
				s.SqlMock.
					ExpectQuery(`SELECT .* FROM "user_forms" WHERE .*`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			nil,
			nil,
			nil,
			&models.PageInfo{},
			[]UserForm{},
		},
		{
			"Find error",
			func() {
				s.SqlMock.
					ExpectQuery(`SELECT .* FROM "user_forms" WHERE .*`).
					WillReturnError(sql.ErrConnDone)
			},
			nil,
			nil,
			sql.ErrConnDone,
			nil,
			nil,
		},
		{
			"Pagination error",
			func() {},
			&negfirst,
			nil,
			ports.ErrInvalidFirst,
			nil,
			nil,
		},
	}

	for _, c := range cases {
		c.Mock()

		conn, err := s.database.UserForm().Pagination(c.First, c.After, nil)
		if c.Error != nil {
			require.ErrorIs(s.T(), err, c.Error, "[%s]", c.Name)
			continue
		}

		require.NoError(s.T(), err, "[%s]", c.Name)
		require.Equal(s.T(), c.PageInfo, conn.PageInfo, "[%s]", c.Name)
		require.Equal(s.T(), len(c.Forms), len(conn.Edges), "[%s]", c.Name)
		for i, edge := range conn.Edges {
			f := c.Forms[i].into()

			assert.NotNil(s.T(), edge, "[%s]", c.Name)
			assert.Equal(s.T(), edge.Cursor, f.ID, "[%s]", c.Name)
			assert.Equal(s.T(), edge.Node, &f, "[%s]", c.Name)
		}
	}
}
