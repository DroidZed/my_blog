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

	fmt.Printf("MONGO URL = %s\n", uri)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB")

	return client
}
