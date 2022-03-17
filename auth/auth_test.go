package auth

import (
	"auction-back/db"
	"auction-back/jwt"
	"auction-back/test"
	"database/sql"
	"net/http"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	os.Setenv("SIGNING_KEY", "test")
	jwt.Init()
}

type AuthSuite struct {
	suite.Suite
	db      *sql.DB
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	handler gin.HandlerFunc
}

func (s *AuthSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.db)
	require.NotNil(s.T(), s.mock)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 s.db,
		PreferSimpleProtocol: true,
	})

	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.handler = New(s.DB)
}

func (s *AuthSuite) TearDownTest() {
	s.db.Close()
}

// Test User in context if token passed
func (s *AuthSuite) TestUser() {
	id := "test-user"
	token, err := jwt.NewUser(id)
	require.NoError(s.T(), err)

	ctx := gin.Context{
		Request: &http.Request{
			Header: http.Header{
				"Authorization": []string{token},
			},
		},
	}

	s.mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE id =").
		WithArgs(id).
		WillReturnRows(test.MockRows(db.User{ID: id}))

	s.handler(&ctx)

	user := ForViewer(ctx.Request.Context())
	require.NotNil(s.T(), user)
	require.Equal(s.T(), user.ID, id)
}

// Test User in context is nil if token passed but id is unknown
func (s *AuthSuite) TestUnknownUser() {
	id := "test-unknown"
	token, err := jwt.NewUser(id)
	require.NoError(s.T(), err)

	ctx := gin.Context{
		Request: &http.Request{
			Header: http.Header{
				"Authorization": []string{token},
			},
		},
	}

	s.mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE id =").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	s.handler(&ctx)

	user := ForViewer(ctx.Request.Context())
	require.Nil(s.T(), user)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
