package user

import (
	"context"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetInput struct {
	filter     bson.M
	projection any
}

type UserProvider interface {
	Add(ctx context.Context, e *User) error
	GetByIdProj(ctx context.Context, id string, proj interface{}) (*User, error)
	GetById(ctx context.Context, id string) (*User, error)
	GetOne(ctx context.Context, in GetInput) (*User, error)
}

type Service struct {
	hasher cryptor.CryptoHelper
	db     *mongo.Database
}

func New(hasher cryptor.CryptoHelper, db *mongo.Database) *Service {
	return &Service{
		hasher: hasher,
		db:     db,
	}
}

func (s *Service) Add(ctx context.Context, entity *User) error {

	modified, err := s.hasher.HashPlain(entity.Password)
	if err != nil {
		return err
	}

	entity.Password = modified

	coll := s.db.Collection("users")

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
			filter: bson.M{"_id": objectId},
			projection: &options.FindOneOptions{
				Projection: projection,
			},
		})
}

func (s *Service) GetById(ctx context.Context, id string) (*User, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, GetInput{filter: bson.M{"_id": objectId}})
}

func (s *Service) GetOne(ctx context.Context, input GetInput) (*User, error) {

	coll := s.db.Collection("users")

	result := &User{}

	if err := coll.FindOne(ctx,
		input.filter,
		&options.FindOneOptions{
			Projection: input.projection,
		},
	).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
