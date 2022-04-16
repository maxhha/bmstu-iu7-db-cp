package server

import (
	"auction-back/adapters/bank"
	"auction-back/adapters/database"
	"auction-back/adapters/role"
	"auction-back/adapters/token"
	"auction-back/adapters/token_mock"
	"auction-back/auth"
	"auction-back/graph"
	"auction-back/graph/generated"
	"auction-back/jwt"
	"auction-back/ports"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var tokenCreators = map[string]func(db ports.DB) ports.Token{
	"PROD": func(db ports.DB) ports.Token {
		senders := []ports.TokenSender{
			emailTokenSender(db),
			phoneTokenSender(db),
		}

		tokenAdapter := token.New(db, senders)
		return &tokenAdapter
	},
	"MOCK": func(db ports.DB) ports.Token {
		tokenAdapter := token_mock.New(db)
		return &tokenAdapter
	},
}

func Init() *gin.Engine {
	filename, ok := os.LookupEnv("SERVER_DOTENV")
	if !ok {
		log.Fatalln("SERVER_DOTENV does not exist in environment variables!")
	}

	if err := godotenv.Load(filename); err != nil {
		log.Fatalln(err)
	}

	jwt.Init()
	db := database.Connect()

	// TODO: Add dependency injection lib
	tokenAdapterName, ok := os.LookupEnv("TOKEN_ADAPTER")
	if !ok || tokenAdapterName == "" {
		tokenAdapterName = "PROD"
	}

	tokenAdapterCreator, exists := tokenCreators[tokenAdapterName]
	if !exists {
		log.Fatalf("Unknown TOKEN_ADAPTER = '%s'\n", tokenAdapterName)
	}

	tokenAdapter := tokenAdapterCreator(&db)
	bankAdapter := bank.New(&db)
	roleAdapter := role.New(&db)

	resolver := graph.New(&db, tokenAdapter, &bankAdapter, &roleAdapter)
	config := generated.Config{Resolvers: resolver}
	config.Directives.HasRole = roleAdapter.Handler()
	schema := generated.NewExecutableSchema(config)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:8080", "http://127.0.0.1:3000", "http://[::1]:3000", "http://[::1]:8080", "http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"POST, GET, OPTIONS"}
	r.Use(cors.New(corsConfig))

	r.Use(auth.New(&db))
	r.Any("/graphql", graphqlHandler(schema))
	r.GET("/graphiql", playgroundHandler())

	return r
}
