package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnection(ctx context.Context, dbUri string) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	clientOptions := options.
		Client().
		ApplyURI(dbUri)

	dbConnection, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server.
	err = dbConnection.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return dbConnection, nil
}
