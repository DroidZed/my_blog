package user

import (
	"context"
	"fmt"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserService interface {
	FindUserByID(id string) (*User, error)
	UpdateOneUser(user User) error
	FindUserByEmail(email string) *User
}

type UserService struct{}

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

func (s *UserService) UpdateOneUser(user User) error {

	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(utils.UserCollection)

	ctx, cancel := context.WithTimeout(context.Background(), utils.ContextTimeOut)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	opt := options.Update().SetUpsert(false)

	log := config.GetLogger()

	update := bson.M{
		"$set": bson.M{
			"fullName": user.FullName,
			"age":      user.Age,
		},
	}

	log.Debug(user)

	updateRes, err := coll.UpdateOne(ctx, filter, update, opt)

	if err != nil {
		return err
	}

	if updateRes.ModifiedCount == 0 {
		return fmt.Errorf("0 modifications happened")
	}

	return nil
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
