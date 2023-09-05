package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/DroidZed/go_lance/internal/user"
)

const collectionName = "users"
const timeOut = 1 * time.Minute

func ValidateUser(login *LoginBody) (string, error) {

	data, err := user.FindUserByEmail(login.Email)
	if err != nil {
		return "", err
	}

	pwdIsValid := cryptor.CompareSecureToPlain(data.Password, login.Password)

	if !pwdIsValid {
		return "", fmt.Errorf("invalid credentials")
	}

	return data.ID.String(), nil
}

func SaveUser(data *user.User) (interface{}, error) {

	coll := config.GetConnection().Database(config.EnvDbName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	modified, err := data.HashUserPassword()
	if err != nil {
		return nil, err
	}

	result, err := coll.InsertOne(ctx, modified)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}
