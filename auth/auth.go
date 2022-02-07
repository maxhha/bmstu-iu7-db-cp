package auth

import (
	"auction-back/db"
	"auction-back/jwt"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

var viewerContextKey = &contextKey{"viewer"}

type contextKey struct {
	name string
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		tokens, ok := c.Request.Header["Authorization"]

		if !ok || len(tokens) == 0 {
			return
		}

		id, err := jwt.Parse(tokens[0])

		fmt.Printf("err: %v\n", err)

		if err != nil {
			return
		}

		fmt.Printf("id: %+v\n", id)

		viewer := db.User{}
		db.DB.Take(&viewer, "id = ?", id)

		ctx := context.WithValue(c.Request.Context(), viewerContextKey, &viewer)
		c.Request = c.Request.WithContext(ctx)
	}
}

func ForViewer(ctx context.Context) *db.User {
	viewer, _ := ctx.Value(viewerContextKey).(*db.User)
	return viewer
}
