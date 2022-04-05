package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn, ok := os.LookupEnv("POSTGRES_CONNECTION")
	if !ok {
		log.Fatalln("POSTGRES_CONNECTION does not exist in environment variables!")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to database!")
	}

	return database
}
