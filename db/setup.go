package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn, ok := os.LookupEnv("POSTGRES_CONNECTION")

	if !ok {
		panic("POSTGRES_CONNECTION does not exist in environment variables!")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}
