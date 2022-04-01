package main

import (
	"auction-back/adapters/notifier"
	"auction-back/auth"
	"auction-back/graph"
	"auction-back/graph/generated"
	"auction-back/jwt"
	"auction-back/ports/bank"
	"auction-back/ports/role"
	"auction-back/ports/token"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"auction-back/db"
)

func graphqlHandler(db *gorm.DB, token *token.TokenPort, bank *bank.BankPort, role *role.RolePort) gin.HandlerFunc {
	resolver := graph.New(db, token, bank, role)
	config := generated.Config{Resolvers: resolver}
	config.Directives.HasRole = role.Handler()

	h := handler.New(generated.NewExecutableSchema(config))

	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})

	h.SetQueryCache(lru.New(1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

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

func init() {
	filename, ok := os.LookupEnv("SERVER_DOTENV")

	if !ok {
		log.Fatalln("SERVER_DOTENV does not exist in environment variables!")
	}

	if err := godotenv.Load(filename); err != nil {
		log.Fatalln(err)
	}

	jwt.Init()
}

func Init() *gin.Engine {
	DB := db.ConnectDatabase()
	notifier := notifier.New("EMAIL_NOTIFIER_ADDRESS")

	tokenPort := token.New(DB, notifier)
	bankPort := bank.New(DB)
	rolePort := role.New(DB)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:8080", "http://127.0.0.1:3000", "http://[::1]:3000", "http://[::1]:8080", "http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"POST, GET, OPTIONS"}
	r.Use(cors.New(corsConfig))

	r.Use(auth.New(DB))
	r.Any("/graphql", graphqlHandler(DB, &tokenPort, &bankPort, &rolePort))
	r.GET("/graphiql", playgroundHandler())

	return r
}

func main() {
	Init().Run()
}
