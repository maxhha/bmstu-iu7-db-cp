package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserGetSuite struct {
	DatabaseSuite
}

func TestUserGetSuite(t *testing.T) {
	suite.Run(t, new(UserGetSuite))
}

func (s *UserGetSuite) TestGetExistedID() {
	id := "test-user"

	s.SqlMock.
		ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 LIMIT 1`).
		WithArgs(id).
		WillReturnRows(MockRows(User{ID: id}))

	result, err := s.database.User().Get(id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), result.ID, id)
}

func (s *UserGetSuite) TestGetUnknownID() {
	id := "unknown-user"

	s.SqlMock.
		ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 LIMIT 1`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, err := s.database.User().Get(id)
	assert.ErrorIs(s.T(), err, sql.ErrNoRows)
}

type UserCreateSuite struct {
	DatabaseSuite
}

func TestUserCreateSuite(t *testing.T) {
	suite.Run(t, new(UserCreateSuite))
}

func (s *UserCreateSuite) TestCreateSuccess() {
	user := models.User{}
	obj := User{ID: "test-user", CreatedAt: time.Now()}

	s.SqlMock.
		ExpectQuery(`INSERT INTO "users" .+ RETURNING "id","created_at"`).
		WithArgs(nil, nil).
		WillReturnRows(MockRows(obj))

	assert.NoError(s.T(), s.database.User().Create(&user))
	assert.Equal(s.T(), obj.into(), user)
}

func (s *UserCreateSuite) TestCreateError() {
	user := models.User{}

	s.SqlMock.
		ExpectQuery(`INSERT INTO "users" .+ RETURNING "id","created_at"`).
		WithArgs(nil, nil).
		WillReturnError(sql.ErrNoRows)

	err := s.database.User().Create(&user)
	assert.ErrorIs(s.T(), err, sql.ErrNoRows)
}

type UserPaginationSuite struct {
	DatabaseSuite
}

func TestUserPaginationSuite(t *testing.T) {
	suite.Run(t, new(UserPaginationSuite))
}

func (s *UserPaginationSuite) TestPaginationSuccessMany() {
	objsN := 10
	objs := make([]User, objsN)
	for i := 0; i < objsN; i++ {
		objs[i].ID = fmt.Sprintf("test-user-%d", i)
		objs[i].CreatedAt = time.Now().Add(time.Duration(i) * time.Hour)
	}

	s.SqlMock.
		ExpectQuery(`SELECT \* FROM "users" ORDER BY created_at desc`).
		WillReturnRows(MockRows(objs))

	config := ports.UserPaginationConfig{}

	conn, err := s.database.User().Pagination(config)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), conn.PageInfo, &models.PageInfo{
		StartCursor: &objs[0].ID,
		EndCursor:   &objs[len(objs)-1].ID,
	})

	for i, obj := range objs {
		assert.Equal(s.T(), conn.Edges[i].Cursor, obj.ID)
		assert.NotNil(s.T(), conn.Edges[i].Node)
		assert.Equal(s.T(), *conn.Edges[i].Node, obj.into())
	}
}

type UserLastApprovedUserFormSuite struct {
	DatabaseSuite
}

func TestUserLastApprovedUserFormSuite(t *testing.T) {
	suite.Run(t, new(UserLastApprovedUserFormSuite))
}

func (s *UserLastApprovedUserFormSuite) TestSuccess() {
	user := models.User{ID: "test-user"}
	form := UserForm{ID: "test-user-form"}

	s.SqlMock.
		ExpectQuery(`SELECT \* FROM "user_forms" WHERE user_id = \$1 AND state = \$2 ORDER BY created_at desc LIMIT 1`).
		WithArgs(user.ID, models.UserFormStateApproved).
		WillReturnRows(MockRows(form))

	result, err := s.database.User().LastApprovedUserForm(user)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), result, form.into())
}
