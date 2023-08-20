package db

import (
	"context"
	"fmt"

	"github.com/DroidZed/go_lance/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnection() *mongo.Client {

	uri := config.EnvDbURI()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to %s\n", client.Database(config.EnvDbName()).Name())

	return client
}

var Client = GetConnection()
