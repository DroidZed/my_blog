package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnection(baseCtx context.Context, env *EnvConfig) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(baseCtx, 1*time.Minute)
	defer cancel()

	clientOptions := options.
		Client().
		ApplyURI(env.DBUri)

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
