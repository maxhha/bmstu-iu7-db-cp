package server

import (
	"auction-back/adapters/bank"
	"auction-back/adapters/database"
	"auction-back/adapters/dealer"
	"auction-back/adapters/role"
	"auction-back/adapters/token"
	"auction-back/auth"
	"auction-back/graph"
	"auction-back/graph/generated"
	"auction-back/jwt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	tokenAdapter := token.New(&db)
	bankAdapter := bank.New(&db)
	roleAdapter := role.New(&db)
	deakerAdapter := dealer.New("http://localhost:8082")

	ownerCheckers := newOwnerCheckers()
	roleCheckers := newRoleCheckers(&roleAdapter, ownerCheckers)

	resolver := graph.New(&db, tokenAdapter, &bankAdapter, &roleAdapter, &deakerAdapter)
	config := generated.Config{Resolvers: resolver}
	config.Directives.HasRole = hasRoleDirective(roleCheckers)
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
