package auth

import (
	"auction-back/db"
	"auction-back/jwt"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

var viewerContextKey = &contextKey{"viewer"}
var guestContextKey = &contextKey{"guest"}

type contextKey struct {
	name string
}

func authUser(token string, c *gin.Context) error {
	id, err := jwt.ParseUser(token)
	if err != nil {
		if err.Error() == "subject is not user" {
			return nil
		}
		return err
	}

	viewer := db.User{}
	if err := db.DB.Take(&viewer, "id = ?", id).Error; err != nil {
		return err
	}

	ctx := context.WithValue(c.Request.Context(), viewerContextKey, &viewer)
	c.Request = c.Request.WithContext(ctx)

	return nil
}

func authGuest(token string, c *gin.Context) error {
	id, err := jwt.ParseGuest(token)
	if err != nil {
		if err.Error() == "subject is not guest" {
			return nil
		}
		return err
	}

	guest := db.Guest{}
	if err := db.DB.Take(&guest, "id = ?", id).Error; err != nil {
		return err
	}

	ctx := context.WithValue(c.Request.Context(), guestContextKey, &guest)
	c.Request = c.Request.WithContext(ctx)

	return nil
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		tokens, ok := c.Request.Header["Authorization"]

		if !ok || len(tokens) == 0 {
			return
		}

		token := tokens[0]

		if err := authGuest(token, c); err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		if err := authUser(token, c); err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
	}
}

func ForViewer(ctx context.Context) *db.User {
	viewer, _ := ctx.Value(viewerContextKey).(*db.User)
	return viewer
}

func ForGuest(ctx context.Context) *db.Guest {
	guest, _ := ctx.Value(guestContextKey).(*db.Guest)
	return guest
}
