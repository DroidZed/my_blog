package signup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/pigeon"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/DroidZed/go_lance/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ISignUpService interface {
	SaveUser(data *user.User) (interface{}, error)
	SaveConfirmationCode(data *ConfirmationCode) (interface{}, error)
	FindCodeByEmail(email string) (interface{}, error)
	VerifyEmail(data *user.User) (interface{}, error)
}

type SignUpService struct{}

const timeOut = 1 * time.Minute

func (s *SignUpService) FindCodeByEmail(email string) (*ConfirmationCode, error) {

	env := config.LoadEnv()
	coll := config.GetConnection().Database(env.DBName).Collection("confirmationCodes")

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

func (s *SignUpService) SaveConfirmationCode(email string) (string, error) {
	env := config.LoadEnv()
	coll := config.GetConnection().Database(env.DBName).Collection("confirmationCodes")

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	code := s.GenerateCode(utils.UPPER_CODE_LIMIT)

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "code",
					Value: code,
				},
				{
					Key:   "email",
					Value: email,
				},
				{
					Key:   "createdAt",
					Value: primitive.NewDateTimeFromTime(time.Now()),
				},
				{
					Key:   "expiresAt",
					Value: primitive.NewDateTimeFromTime(time.Now().Add(time.Duration(time.Duration.Minutes(15)))),
				},
			},
		},
	}

	filter := bson.M{"email": email}

	opts := options.Update().SetUpsert(true)

	updateRes, updateErr := coll.UpdateOne(ctx, filter, update, opts)
	if updateErr != nil {
		return "", updateErr
	}

	if updateRes.ModifiedCount == 0 {
		return "", fmt.Errorf("0 modifications happened")
	}

	return code, nil
}

func (s *SignUpService) CheckCodeValidity(email string) bool {

	verifyCode, err := s.FindCodeByEmail(email)

	if err != nil {
		return false
	}

	return time.Now().After(verifyCode.ExpiresAt.Time())
}

func (s *SignUpService) deliverEmailToUser(
	to,
	subject,
	templateName,
	confEntity string,
) error {

	req := pigeon.NewRequest(
		[]string{to},
		subject,
		"",
	)

	err := req.ParseTemplate(templateName, confEntity)

	if err != nil {
		return err
	}

	emailErr := req.SendEmail()

	if emailErr != nil {
		return emailErr
	}

	return nil
}

func (s *SignUpService) GenerateCode(bound int64) string {
	builder := &strings.Builder{}

	for i := 0; i < 4; i++ {
		builder.WriteString(fmt.Sprintf("%d", utils.RNG(bound)))
	}

	return builder.String()
}

func (s *SignUpService) SaveCodeAndSendEmail(

	email string,
) error {

	code, err := s.SaveConfirmationCode(email)

	if err != nil {
		return err
	}

	if emailErr := s.deliverEmailToUser(
		email,
		"CONFIRMATION MAIL",
		"confirmation_email",
		code,
	); emailErr != nil {
		return emailErr
	}

	return nil
}
