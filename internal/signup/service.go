package signup

import (
	"context"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/user"
	"go.mongodb.org/mongo-driver/bson"
)

type ISignUpService interface {
	SaveUser(data *user.User) (interface{}, error)
	SaveConfirmationCode(data *ConfirmationCode) (interface{}, error)
	FindCodeByEmail(email string) (interface{}, error)
	VerifyEmail(data *user.User) (interface{}, error)
}

type SignUpService struct{}

const timeOut = 1 * time.Minute

func (s *SignUpService) SaveUser(data *user.User) (interface{}, error) {

	env := config.LoadEnv()

	coll := config.GetConnection().Database(env.DBName).Collection("users")

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

func (s *SignUpService) DeleteCode(data *user.User) (interface{}, error) {
	return nil, nil
}

func (s *SignUpService) FindCodeByEmail(email string) (*ConfirmationCode, error) {

	env := config.LoadEnv()
	coll := config.GetConnection().Database(env.DBName).Collection("confirmationTokens")

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	result := &ConfirmationCode{}

	filter := bson.M{"email": email}

	err := coll.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SignUpService) SaveConfirmationCode(data *ConfirmationCode) (*ConfirmationCode, error) {
	return nil, nil
}
