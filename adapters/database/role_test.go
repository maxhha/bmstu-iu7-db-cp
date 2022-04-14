package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleSuite struct {
	DatabaseSuite
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

func (s *RoleSuite) TestFind() {
	role1 := Role{Type: models.RoleTypeAdmin}
	role2 := Role{Type: models.RoleTypeManager}
	role3 := Role{UserID: "test-1"}

	cases := []struct {
		Name   string
		Config ports.RoleFindConfig
		Error  error
		Mock   func()
		Roles  []Role
	}{
		{
			Name: "No filter",
			Mock: func() {
				s.SqlMock.
					ExpectQuery(`SELECT \* FROM "roles"`).
					WillReturnRows(MockRows(role1, role2, role3))
			},
			Roles: []Role{role1, role2, role3},
		},
		{
			Name: "With limit",
			Config: ports.RoleFindConfig{
				Limit: 10,
			},
			Mock: func() {
				s.SqlMock.
					ExpectQuery(`SELECT \* FROM "roles" WHERE "roles"\."deleted_at" IS NULL LIMIT 10`).
					WillReturnRows(MockRows(role1, role2, role3))
			},
			Roles: []Role{role1, role2, role3},
		},
		{
			Name: "With types",
			Config: ports.RoleFindConfig{
				Types: []models.RoleType{models.RoleTypeManager},
			},
			Mock: func() {
				s.SqlMock.
					ExpectQuery(`SELECT \* FROM "roles" WHERE type IN \(\$1\)`).
					WithArgs(models.RoleTypeManager).
					WillReturnRows(MockRows(role2))
			},
			Roles: []Role{role2},
		},
		{
			Name: "With userIds",
			Config: ports.RoleFindConfig{
				UserIDs: []string{role3.UserID},
			},
			Mock: func() {
				s.SqlMock.
					ExpectQuery(`SELECT \* FROM "roles" WHERE user_id IN \(\$1\)`).
					WithArgs(role3.UserID).
					WillReturnRows(MockRows(role3))
			},
			Roles: []Role{role3},
		},
		{
			Name:   "Find error",
			Config: ports.RoleFindConfig{},
			Mock: func() {
				s.SqlMock.
					ExpectQuery(`SELECT \* FROM "roles"`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			Error: sql.ErrNoRows,
		},
	}

	for _, c := range cases {
		c.Mock()
		roles, err := s.database.Role().Find(c.Config)

		if c.Error != nil {
			require.ErrorIs(s.T(), err, c.Error, "[%s]", c.Name)
			require.Nil(s.T(), roles, "[%s]", c.Name)
			continue
		}

		require.NoError(s.T(), err)
		for i, role := range roles {
			require.Equal(s.T(), role, c.Roles[i].into(), "[%s]", c.Name)
		}
	}
}
