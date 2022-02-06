package auth

import (
	"auction-back/db"
	"context"

	"github.com/gin-gonic/gin"
)

var viewerContextKey = &contextKey{"viewer"}

type contextKey struct {
	name string
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokens, ok := c.Request.Header["Authorization"]

		if ok && len(tokens) == 0 {
			viewer := db.User{}
			db.DB.Take(&viewer, tokens[0])

			ctx := context.WithValue(c.Request.Context(), viewerContextKey, &viewer)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

func ForViewer(ctx context.Context) *db.User {
	viewer, _ := ctx.Value(viewerContextKey).(*db.User)
	return viewer
}
