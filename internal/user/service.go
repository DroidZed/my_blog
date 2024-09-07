package user

import (
	"context"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"go.mongodb.org/mongo-driver/bson"
)

type UserProvider interface {
	SaveUser(ctx context.Context, data *User) error
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service struct {
	Hasher   cryptor.CryptoHelper
	UserRepo UserRepoProvider
}

func (s *Service) SaveUser(ctx context.Context, data *User) error {

	modified, err := s.Hasher.HashPlain(data.Password)
	if err != nil {
		return err
	}

	data.Password = modified

	return s.UserRepo.Save(ctx, data)
}

func (s *Service) FindUserByID(ctx context.Context, id string) (*User, error) {

	return s.UserRepo.FindById(ctx, id)
}

func (s *Service) FindUserByEmail(ctx context.Context, email string) (*User, error) {

	return s.UserRepo.FindOne(ctx, bson.M{"email": email}, nil)
}
