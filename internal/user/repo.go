package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepoProvider interface {
	Save(ctx context.Context, entity *User) error
	FindOne(ctx context.Context, filter bson.M, projection interface{}) (*User, error)
	FindById(ctx context.Context, id string) (*User, error)
}

type UserRepo struct {
	DbClient *mongo.Client
	DBName   string
}

func (r *UserRepo) Save(ctx context.Context, e *User) error {

	coll := r.DbClient.Database(r.DBName).Collection("users")

	_, insertErr := coll.InsertOne(ctx, e)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (r *UserRepo) FindOne(ctx context.Context, filter bson.M, projection interface{}) (*User, error) {

	coll := r.DbClient.Database(r.DBName).Collection("users")

	result := &User{}

	opts := &options.FindOneOptions{
		Projection: projection,
	}

	if err := coll.FindOne(ctx, filter, opts).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserRepo) FindById(ctx context.Context, id string) (*User, error) {
	coll := r.DbClient.Database(r.DBName).Collection("users")

	result := &User{}

	objectId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return nil, err1
	}

	if err := coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
