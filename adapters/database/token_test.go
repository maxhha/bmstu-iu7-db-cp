package database

import (
	"auction-back/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type TokenCreateSuite struct {
	DatabaseSuite
}

func TestTokenCreateSuite(t *testing.T) {
	suite.Run(t, new(TokenCreateSuite))
}

func (s *TokenCreateSuite) TestCreateSuccess() {
	token := models.Token{
		UserID:    "test-user",
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(3)),
		Action:    models.TokenActionSetUserEmail,
		Data:      map[string]interface{}{"hello": "world"},
	}
	obj := Token{
		ID:        31337,
		UserID:    token.UserID,
		CreatedAt: time.Now(),
		ExpiresAt: token.ExpiresAt,
		Action:    token.Action,
		Data:      token.Data,
	}

	s.SqlMock.
		ExpectQuery(`INSERT INTO "tokens" .+ RETURNING "created_at","id"`).
		WithArgs(token.UserID, nil, token.ExpiresAt, token.Action, datatypes.JSONMap(token.Data)).
		WillReturnRows(MockRows(obj))

	assert.NoError(s.T(), s.database.Token().Create(&token))
	assert.Equal(s.T(), obj.into(), token)
}
