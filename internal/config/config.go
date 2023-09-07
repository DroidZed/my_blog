package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port          int64
	Host          string
	DBUri         string
	DBName        string
	AccessSecret  string
	AccessExpiry  string
	RefreshSecret string
	RefreshExpiry string
}

var config *EnvConfig

func LoadConfig() *EnvConfig {

	if config != nil {
		return config
	}

	log := InitializeLogger().LogHandler

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading work environment.\n %s", err)
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		log.Fatalf("Error loading work environment.\n %s", err)

	}

	conf := &EnvConfig{
		Port:          port,
		Host:          os.Getenv("HOST"),
		DBUri:         os.Getenv("DB_URI"),
		DBName:        os.Getenv("DB_NAME"),
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		AccessExpiry:  os.Getenv("ACCESS_EXPIRY"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		RefreshExpiry: os.Getenv("REFRESH_EXPIRY"),
	}

	config = conf

	return config
}
