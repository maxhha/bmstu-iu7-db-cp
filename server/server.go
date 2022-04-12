package server

import (
	"auction-back/adapters/bank"
	"auction-back/adapters/database"
	"auction-back/adapters/role"
	"auction-back/adapters/token"
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
	db := database.Connect()

	senders := []ports.TokenSender{
		emailTokenSender(),
		phoneTokenSender(&db),
	}

	tokenPort := token.New(&db, senders)
	bankPort := bank.New(&db)
	rolePort := role.New(&db)

	resolver := graph.New(&db, &tokenPort, &bankPort, &rolePort)
	config := generated.Config{Resolvers: resolver}
	config.Directives.HasRole = rolePort.Handler()

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:8080", "http://127.0.0.1:3000", "http://[::1]:3000", "http://[::1]:8080", "http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"POST, GET, OPTIONS"}
	r.Use(cors.New(corsConfig))

	r.Use(auth.New(&db))
	r.Any("/graphql", graphqlHandler(config))
	r.GET("/graphiql", playgroundHandler())

	return r
}
