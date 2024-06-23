package user

import (
	"context"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserService interface {
	SaveUser(data *User) error
	FindUserByID(id string) (*User, error)
	FindUserByEmail(email string) *User
}

type UserService struct{}

func (s *UserService) SaveUser(data *User) error {

	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(utils.UserCollection)

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

func (s *UserService) FindUserByID(id string) (*User, error) {
	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(utils.UserCollection)

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

func (s *UserService) FindUserByEmail(email string) *User {
	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(utils.UserCollection)

	ctx, cancel := context.WithTimeout(context.Background(), utils.ContextTimeOut)
	defer cancel()

	result := &User{}

	filter := bson.M{"email": email}

	if err := coll.FindOne(ctx, filter).Decode(result); err != nil {
		return nil
	}

	return result
}
