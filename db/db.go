package db

import (
	"context"
	"time"

	"github.com/DroidZed/go_lance/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnection() *mongo.Client {

	log := config.Logger.LogHandler

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.EnvDbURI()))
	if err != nil {
		log.Fatal("Could not connect to database.")
	}

	// Ping the MongoDB server.
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Infof("Connected to %s\n", config.EnvDbName())

	return client
}

var Client = GetConnection()
