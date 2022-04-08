package auth

import (
	"auction-back/jwt"
	"auction-back/models"
	"auction-back/ports"
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var ErrUnauthorized = errors.New("unauthorized")

type auth struct {
	db ports.DB
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

	viewer, err := a.db.User().Get(id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	ctx := WithViewer(c.Request.Context(), viewer)
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

func New(db ports.DB) gin.HandlerFunc {
	a := auth{db}
	return func(c *gin.Context) {
		a.apply(c)
	}
}

func ForViewer(ctx context.Context) (models.User, error) {
	viewer, ok := ctx.Value(viewerContextKey).(models.User)
	if !ok {
		return models.User{}, ErrUnauthorized
	}
	return viewer, nil
}

func WithViewer(c context.Context, viewer models.User) context.Context {
	return context.WithValue(c, viewerContextKey, viewer)
}
