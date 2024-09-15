package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port          int64
	Env           string
	Host          string
	DBUri         string
	DBName        string
	AccessSecret  string
	AccessExpiry  string
	RefreshSecret string
	RefreshExpiry string
	SmtpHost      string
	SmtpPort      string
	SmtpUsername  string
	SmtpPassword  string
	MASTER_PWD    string
	MASTER_EMAIL  string
}

func LoadEnv() (*EnvConfig, error) {

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)

	if err != nil {
		return nil, err
	}

	if err := loadVarsPerEnv(); err != nil {
		return nil, err

	}

	config := &EnvConfig{
		Port:          port,
		DBUri:         os.Getenv("DB_URI"),
		Env:           os.Getenv("ENV"),
		Host:          os.Getenv("HOST"),
		DBName:        os.Getenv("DB_NAME"),
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		AccessExpiry:  os.Getenv("ACCESS_EXPIRY"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		RefreshExpiry: os.Getenv("REFRESH_EXPIRY"),
		SmtpHost:      os.Getenv("SMTP_HOST"),
		SmtpPort:      os.Getenv("SMTP_PORT"),
		SmtpUsername:  os.Getenv("SMTP_USERNAME"),
		SmtpPassword:  os.Getenv("SMTP_PASSWORD"),
		MASTER_PWD:    os.Getenv("MASTER_PWD"),
		MASTER_EMAIL:  os.Getenv("MASTER_EMAIL"),
	}

	return config, nil

}

func loadVarsPerEnv() (err error) {

	env, ok := os.LookupEnv("ENV")

	if !ok {
		return fmt.Errorf("environment not set")
	}

	switch env {
	case "dev":

		if err := godotenv.Load(".env.dev"); err != nil {
			return err
		}

	case "prod":

		if err := godotenv.Load(".env.prod"); err != nil {
			return err
		}

	default:

		return fmt.Errorf("environment not set")

	}

	return nil
}
