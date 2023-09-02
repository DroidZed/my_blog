package user

import (
	"context"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllUsers() []User {
	log := config.InitializeLogger().LogHandler

	coll := db.GetConnection().Database(config.EnvDbName()).Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := coll.Find(ctx, bson.D{})

	results := make([]User, 0)

	if err != nil {
		return nil
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err2 := cur.Close(ctx)
		if err2 != nil {
			log.Errorf("Error closing cursor.\n%s", err2.Error())
		}
	}(cur, ctx)

	for {
		if cur.TryNext(ctx) {
			var doc User
			err := cur.Decode(&doc)
			if err != nil {
				return nil
			}
			results = append(results, doc)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		if cur.ID() == 0 {
			break
		}
	}

	return results
}

func FindUserById(id string) *User {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection("users")

	filter := bson.D{{Key: "_id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := &User{}

	err := coll.FindOne(ctx, filter).Decode(*result)

	if err != nil {
		return nil
	}

	return result
}

func SaveOne(data *User) interface{} {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, data)

	if err != nil {
		return nil
	}

	return result.InsertedID

}

func DeleteOne(id string) bool {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection("users")

	filter := bson.D{{Key: "_id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := coll.DeleteOne(ctx, filter)

	if err != nil || result.DeletedCount == 0 {
		return false
	}

	return true
}

func UpdateOne(id string, user *User) error {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := coll.UpdateByID(ctx, id, user, options.Update().SetUpsert(false))

	if result.ModifiedCount == 0 || err != nil {
		return err
	}

	return nil
}
