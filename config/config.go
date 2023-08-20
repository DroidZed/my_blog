package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvDbURI() string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DB_URI")
}

func EnvDbPORT() (int64, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	return strconv.ParseInt(os.Getenv("PORT"), 10, 64)
}
