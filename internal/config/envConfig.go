package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port          int64
	DBUri         string
	DBName        string
	AccessSecret  string
	AccessExpiry  string
	RefreshSecret string
	RefreshExpiry string
	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_USERNAME string
	SMTP_PASSWORD string
}

var config *EnvConfig

func LoadEnv() *EnvConfig {

	if config != nil {
		return config
	}

	log := InitializeLogger().LogHandler

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading work environment.\n %s", err)
	}

	if port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64); err != nil {
		log.Fatalf("Error loading work environment.\n %s", err)

	} else {

		if err := loadVarsPerEnv(); err != nil {
			log.Fatalf("Error loading work environment.\n %s", err)
		}

		conf := &EnvConfig{
			Port:          port,
			DBUri:         os.Getenv("DB_URI"),
			DBName:        os.Getenv("DB_NAME"),
			AccessSecret:  os.Getenv("ACCESS_SECRET"),
			AccessExpiry:  os.Getenv("ACCESS_EXPIRY"),
			RefreshSecret: os.Getenv("REFRESH_SECRET"),
			RefreshExpiry: os.Getenv("REFRESH_EXPIRY"),
			SMTP_HOST:     os.Getenv("SMTP_HOST"),
			SMTP_PORT:     os.Getenv("SMTP_PORT"),
			SMTP_USERNAME: os.Getenv("SMTP_USERNAME"),
			SMTP_PASSWORD: os.Getenv("SMTP_PASSWORD"),
		}

		config = conf

	}
	return config
}

func loadVarsPerEnv() error {

	if env, ok := os.LookupEnv("ENV"); !ok {
		return fmt.Errorf("environment not set")
	} else {

		switch env {
		case "dev":
			{
				if err := godotenv.Load(".env.dev"); err != nil {
					return err
				}
			}
		case "prod":
			{
				if err := godotenv.Load(".env.prod"); err != nil {
					return err
				}
			}
		default:
			{
				return fmt.Errorf("environment not set")
			}
		}
	}
	return nil
}
