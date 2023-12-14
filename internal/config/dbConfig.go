package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var globalDbHandle *mongo.Client

func GetConnection() *mongo.Client {

	if globalDbHandle != nil {
		return globalDbHandle
	}

	log := InitializeLogger().LogHandler

	env := LoadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dbConnection, err := mongo.Connect(ctx, options.Client().ApplyURI(env.DBUri))
	if err != nil {
		log.Fatal("Could not connect to database.")
	}

	// Ping the MongoDB server.
	err = dbConnection.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Infof("Connected to %s\n", env.DBName)

	globalDbHandle = dbConnection

	return globalDbHandle
}
