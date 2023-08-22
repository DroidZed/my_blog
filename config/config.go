package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const EnvFile = ".env"

func _LoadEnv() {
	if err := godotenv.Load(EnvFile); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func EnvDbURI() string {
	_LoadEnv()
	return os.Getenv("DB_URI")
}

func EnvDbPORT() (int64, error) {
	_LoadEnv()
	return strconv.ParseInt(os.Getenv("PORT"), 10, 64)
}

func EnvDbName() string {
	_LoadEnv()
	return os.Getenv("DB_NAME")
}

func EnvHost() string {
	_LoadEnv()
	return os.Getenv("HOST")
}
