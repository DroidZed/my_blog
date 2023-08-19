package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvDbURI() string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DB_URI")
}

func EnvDbPORT() string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("PORT")
}
