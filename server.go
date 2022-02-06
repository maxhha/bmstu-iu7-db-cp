package main

import (
	"auction-back/auth"
	"auction-back/graph"
	"auction-back/graph/generated"

	"github.com/gin-gonic/gin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"auction-back/db"
)

func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	db.ConnectDatabase()

	r := gin.Default()
	r.Use(auth.Middleware())
	r.POST("/graphql", graphqlHandler())
	r.GET("/graphql", playgroundHandler())
	r.Run()
}
