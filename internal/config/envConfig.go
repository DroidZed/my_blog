package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port   int64
	Env    string
	Host   string
	DBUri  string
	DBName string
}

func LoadEnv() (*EnvConfig, error) {

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)

	if err != nil {
		return nil, err
	}

	config := &EnvConfig{
		Port:   port,
		DBUri:  os.Getenv("DB_URI"),
		Env:    os.Getenv("APP_ENV"),
		Host:   os.Getenv("HOST"),
		DBName: os.Getenv("DB_NAME"),
	}

	return config, nil
}
