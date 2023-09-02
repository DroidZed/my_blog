package db

import (
	"context"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbConnection *mongo.Client

func GetConnection() *mongo.Client {

	log := config.InitializeLogger().LogHandler

	if dbConnection != nil {
		return dbConnection
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbConnection, err := mongo.Connect(ctx, options.Client().ApplyURI(config.EnvDbURI()))
	if err != nil {
		log.Fatal("Could not connect to database.")
	}

	// Ping the MongoDB server.
	err = dbConnection.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Infof("Connected to %s\n", config.EnvDbName())

	return dbConnection
}
