package main

import (
	"auction-back/auth"
	"auction-back/jwt"
	"auction-back/ports/bank"
	"auction-back/ports/role"
	"auction-back/ports/token"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"auction-back/db"
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
	DB := db.ConnectDatabase()

	senders := []token.SenderInterface{
		emailTokenSender(),
		phoneTokenSender(),
	}

	tokenPort := token.New(DB, senders)
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
