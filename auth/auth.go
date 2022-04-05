package auth

import (
	"auction-back/jwt"
	"auction-back/models"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type auth struct {
	db *gorm.DB
}

type contextKey struct {
	name string
}

var viewerContextKey = &contextKey{"viewer"}

func (a *auth) authUser(token string, c *gin.Context) error {
	id, err := jwt.ParseUser(token)
	if err != nil {
		return err
	}

	viewer := models.User{}
	if err := a.db.Take(&viewer, "id = ?", id).Error; err != nil {
		return err
	}

	ctx := WithViewer(c.Request.Context(), &viewer)
	c.Request = c.Request.WithContext(ctx)

	return nil
}

func (a *auth) apply(c *gin.Context) {
	defer c.Next()

	tokens, ok := c.Request.Header["Authorization"]

	if !ok || len(tokens) == 0 {
		return
	}

	token := tokens[0]

	if err := a.authUser(token, c); err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "auth err: %v\n", err)
		return
	}
}

func New(db *gorm.DB) gin.HandlerFunc {
	a := auth{db}
	return func(c *gin.Context) {
		a.apply(c)
	}
}

func ForViewer(ctx context.Context) *models.User {
	viewer, _ := ctx.Value(viewerContextKey).(*models.User)
	return viewer
}

func WithViewer(c context.Context, viewer *models.User) context.Context {
	return context.WithValue(c, viewerContextKey, viewer)
}
