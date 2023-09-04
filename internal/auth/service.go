package auth

import (
	"context"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/db"
	"github.com/DroidZed/go_lance/internal/user"
)

const collectionName = "users"
const timeOut = 1 * time.Minute

func SaveUser(data *user.User) (interface{}, error) {

	coll := db.GetConnection().Database(config.EnvDbName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	result, err := coll.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}
