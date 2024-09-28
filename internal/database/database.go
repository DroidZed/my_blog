package database

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	name = os.Getenv("DB_NAME")
)

type DatabaseManager interface {
	Health(ctx context.Context) (bool, error)
}

type Service struct {
	Client *mongo.Client
	Name   string
}

func New(ctx context.Context) (*Service, error) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s/%s", host, port, name)

	dbConnection, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &Service{
		Client: dbConnection,
		Name:   name,
	}, nil

}

func (s *Service) Health(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := s.Client.Ping(ctx, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}
