package user

import (
	"context"
	"fmt"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	hasher cryptor.CryptoHelper
	client *mongo.Client
	dbName string
}

func NewService(hasher cryptor.CryptoHelper, db *mongo.Client, dbName string) *Service {
	return &Service{
		hasher: hasher,
		client: db,
		dbName: dbName,
	}
}

func (s *Service) Add(ctx context.Context, entity User) error {

	coll := s.client.Database(s.dbName).Collection("users")

	_, insertErr := coll.InsertOne(ctx, entity)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (s *Service) GetByIdProj(
	ctx context.Context,
	id string,
	projection any,
) (*User, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.GetOne(
		ctx,
		GetInput{
			Filter: bson.M{"_id": objectId},
			Projection: &options.FindOneOptions{
				Projection: projection,
			},
		})
}

func (s *Service) GetById(ctx context.Context, id string) (*User, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, GetInput{Filter: bson.M{"_id": objectId}})
}

func (s *Service) GetOne(ctx context.Context, input GetInput) (*User, error) {

	coll := s.client.Database(s.dbName).Collection("users")

	result := &User{}

	if err := coll.FindOne(ctx,
		input.Filter,
		&options.FindOneOptions{
			Projection: input.Projection,
		},
	).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) Validate(ctx context.Context, email, password string) (*User, error) {

	data, err := s.GetOne(
		ctx,
		GetInput{
			Filter: bson.M{"email": email},
		},
	)

	if err != nil {
		return nil, err
	}

	pwdIsValid := s.hasher.CompareSecureToPlain(data.Password, password)

	if !pwdIsValid {
		return nil, fmt.Errorf("invalid credentials")
	}

	return data, nil
}
