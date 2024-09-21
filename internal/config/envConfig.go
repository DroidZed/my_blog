package config

import (
	"os"
	"strconv"
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

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)

	if err != nil {
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
