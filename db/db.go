package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DroidZed/go_lance/config"
	"github.com/withmandala/go-log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logger = log.New(os.Stderr)

func GetConnection() *mongo.Client {

	uri := config.EnvDbURI()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Fatal("Could not connect to database.")

	}

	// Ping the MongoDB server.
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal("Failed to ping the database:", err)

	}

	fmt.Printf("Connected to %s\n", client.Database(config.EnvDbName()).Name())

	return client
}

var Client = GetConnection()
