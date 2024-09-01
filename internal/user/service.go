package user

import (
	"context"

	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserProvider interface {
	SaveUser(data *User) error
	FindUserByID(id string) (*User, error)
	FindUserByEmail(email string) *User
}

type Service struct {
	DbClient *mongo.Client
	DBName   string
}

func (s *Service) SaveUser(data *User) error {

	coll := s.DbClient.Database(s.DBName).Collection(utils.UserCollection)

	ctx, cancel := context.WithTimeout(context.Background(), utils.ContextTimeOut)
	defer cancel()

	modified, err := data.HashUserPassword()
	if err != nil {
		return err
	}

	_, insertErr := coll.InsertOne(ctx, modified)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (s *Service) FindUserByID(id string) (*User, error) {

	coll := s.DbClient.Database(s.DBName).Collection(utils.UserCollection)

	ctx, cancel := context.WithTimeout(context.Background(), utils.ContextTimeOut)
	defer cancel()

	result := &User{}

	objectId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return nil, err1
	}

	filter := bson.M{"_id": objectId}

	// Check for errors when executing FindOne
	err := coll.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) FindUserByEmail(email string) *User {

	coll := s.DbClient.Database(s.DBName).Collection(utils.UserCollection)

	ctx, cancel := context.WithTimeout(context.Background(), utils.ContextTimeOut)
	defer cancel()

	result := &User{}

	filter := bson.M{"email": email}

	if err := coll.FindOne(ctx, filter).Decode(result); err != nil {
		return nil
	}

	return result
}
