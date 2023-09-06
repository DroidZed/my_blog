package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const EnvFile = ".env"

func LoadEnv() {
	if err := godotenv.Load(EnvFile); err != nil {
		log.Fatalf("Error loading work environment.\n %s", err)
	}
}

func EnvDbURI() string {
	return os.Getenv("DB_URI")
}

func EnvDbPORT() (int64, error) {
	return strconv.ParseInt(os.Getenv("PORT"), 10, 64)
}

func EnvDbName() string {
	return os.Getenv("DB_NAME")
}

func EnvHost() string {
	return os.Getenv("HOST")
}

func EnvJwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

func EnvJwtExp() string {
	return os.Getenv("JWT_EXPIRY")
}

func EnvRefreshSecret() string {
	return os.Getenv("REFRESH_SECRET")
}

func EnvRefreshExp() string {
	return os.Getenv("REFRESH_EXPIRY")
}
