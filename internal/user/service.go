package user

import (
	"context"
	"fmt"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "users"
const timeOut = 1 * time.Minute

type IUserService interface {
	FindAllUsers() ([]User, error)
	FindUserByID(id string) (*User, error)
	FindUserByEmail(email string) (*User, error)
	UpdateOneUser(user User) error
	DeleteOne(id string) bool
}

type UserService struct{}

func (s *UserService) SaveUser(data *User) error {

	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
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

func (s *UserService) FindAllUsers() ([]User, error) {
	env := config.LoadEnv()

	log := config.InitializeLogger().LogHandler

	coll := config.GetConnection().Database(env.DBName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer cur.Close(ctx)

	results := make([]User, 0)

	for cur.Next(ctx) {
		doc := &User{}
		err := cur.Decode(&doc)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		results = append(results, *doc)
	}

	if err := cur.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	return results, nil
}

func (s *UserService) FindUserByID(id string) (*User, error) {
	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
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

	coll := config.GetConnection().Database(env.DBName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	result := &User{}

	filter := bson.M{"email": email}

	if err := coll.FindOne(ctx, filter).Decode(result); err != nil {
		return nil
	}

	return result
}

func (s *UserService) UpdateOneUser(user User) error {

	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	opt := options.Update().SetUpsert(false)

	log := config.InitializeLogger().LogHandler

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

func (s *UserService) DeleteOne(id string) error {
	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectId}

	result, err := coll.DeleteOne(ctx, filter)

	if err != nil || result.DeletedCount == 0 {
		return err
	}

	return nil
}

func (s *UserService) ActivateUserAccount(email string) error {

	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	filter := bson.M{"email": email}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "accStatus",
					Value: true,
				},
			},
		},
	}

	updateRes, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if updateRes.ModifiedCount == 0 {
		return fmt.Errorf("0 modifications happened")
	}

	return nil
}
