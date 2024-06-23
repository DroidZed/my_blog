package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnection() *mongo.Client {

	log := GetLogger()

	env := LoadEnv()

	log.Info("Opening a database connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	clientOptions := options.
		Client().
		ApplyURI(env.DBUri)

	dbConnection, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Could not connect to the database !")
	}

	// Ping the MongoDB server.
	err = dbConnection.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Infof("Connection to %s has been established.\n", env.DBName)

	connectedHandle := dbConnection

	return connectedHandle
}
