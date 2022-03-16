package auth

import (
	"auction-back/db"
	"auction-back/jwt"
	"auction-back/test"
	"database/sql"
	"net/http"
	"os"
	"testing"
	"time"

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
	db   *sql.DB
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	a    Auth
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

	s.a = New(s.DB)
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

	s.a.Middleware()(&ctx)

	user := ForViewer(ctx.Request.Context())
	require.NotNil(s.T(), user)
	require.Equal(s.T(), user.ID, id)

	guest := ForGuest(ctx.Request.Context())
	require.Nil(s.T(), guest)
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

	s.a.Middleware()(&ctx)

	user := ForViewer(ctx.Request.Context())
	require.Nil(s.T(), user)
}

// Test Guest in context if token passed
func (s *AuthSuite) TestGuest() {
	id := "test-guest"
	token, err := jwt.NewGuest(id, time.Now().Add(time.Hour*time.Duration(1)))
	require.NoError(s.T(), err)

	ctx := gin.Context{
		Request: &http.Request{
			Header: http.Header{
				"Authorization": []string{token},
			},
		},
	}

	s.mock.ExpectQuery("SELECT \\* FROM \"guests\" WHERE id =").
		WithArgs(id).
		WillReturnRows(test.MockRows(db.Guest{ID: id}))

	s.a.Middleware()(&ctx)

	guest := ForGuest(ctx.Request.Context())
	require.NotNil(s.T(), guest)
	require.Equal(s.T(), guest.ID, id)

	user := ForViewer(ctx.Request.Context())
	require.Nil(s.T(), user)
}

// Test Guest in context is nil if token passed but id is unknown
func (s *AuthSuite) TestUnknownGuest() {
	id := "test-unknown"
	token, err := jwt.NewGuest(id, time.Now().Add(time.Hour*time.Duration(1)))
	require.NoError(s.T(), err)

	ctx := gin.Context{
		Request: &http.Request{
			Header: http.Header{
				"Authorization": []string{token},
			},
		},
	}

	s.mock.ExpectQuery("SELECT \\* FROM \"guest\" WHERE id =").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	s.a.Middleware()(&ctx)

	guest := ForGuest(ctx.Request.Context())
	require.Nil(s.T(), guest)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
