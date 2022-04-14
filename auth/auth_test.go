package auth

import (
	"auction-back/jwt"
	"auction-back/models"
	"auction-back/ports"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func init() {
	os.Setenv("SIGNING_KEY", "test")
	jwt.Init()
}

type AuthSuite struct {
	suite.Suite
	db      ports.DBMock
	handler gin.HandlerFunc
}

func (s *AuthSuite) SetupTest() {
	s.handler = New(&s.db)
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

	s.db.UserMock.On("Get", id).Return(models.User{ID: id}, nil)

	s.handler(&ctx)

	user, err := ForViewer(ctx.Request.Context())
	require.NoError(s.T(), err)
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

	s.db.UserMock.On("Get", id).Return(models.User{}, ports.ErrRecordNotFound)

	s.handler(&ctx)

	_, err = ForViewer(ctx.Request.Context())
	require.ErrorIs(s.T(), err, ErrUnauthorized)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
