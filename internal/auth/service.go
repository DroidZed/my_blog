package auth

import (
	"context"
	"log/slog"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userProvider interface {
	Add(ctx context.Context, u user.User) error
	GetByIDProj(ctx context.Context, id string, in utils.GetInput) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetOne(ctx context.Context, in utils.GetInput) (*user.User, error)
	Validate(ctx context.Context, email, password string) (*user.User, error)
}

type Service struct {
	userSrv    userProvider
	refreshKey string
	hasher     cryptor.CryptoHelper
	logger     *slog.Logger

	contactEmail string
	password     string
}

func NewService(
	userSrv userProvider,
	refreshKey string,
	hasher cryptor.CryptoHelper,
	logger *slog.Logger,
	contactEmail string,
	password string,
) *Service {
	return &Service{
		userSrv:    userSrv,
		refreshKey: refreshKey,
		hasher:     hasher,
		logger:     logger,

		contactEmail: contactEmail,
		password:     password,
	}
}

func (s *Service) GenerateNewTokens(expiredToken string) (string, string, error) {
	access, err := s.hasher.ParseToken(expiredToken, s.refreshKey)
	if err != nil {
		return "", "", err
	}

	userID, err := s.hasher.ExtractSubFromClaims(access)
	if err != nil {
		return "", "", err
	}

	newAcc, err := s.hasher.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	newRef, err1 := s.hasher.GenerateRefreshToken()
	if err1 != nil {
		return "", "", err
	}

	return newAcc, newRef, nil
}

func (s *Service) CreateLoginResponse(ctx context.Context, body LoginBody) (LoginResponse, error) {

	user, err := s.userSrv.Validate(ctx, body.Email, body.Password)

	errLogin := LoginResponse{Error: "Unable to login, please try again later"}

	if err != nil {
		return errLogin, err
	}

	userID := user.ID.Hex()

	access, err := s.hasher.GenerateAccessToken(userID)
	if err != nil {
		return errLogin, err
	}

	refresh, err1 := s.hasher.GenerateRefreshToken()
	if err1 != nil {
		return errLogin, err
	}

	if access == "" || refresh == "" {
		return errLogin, err
	}

	return LoginResponse{Jwt: access, Refresh: refresh}, nil
}

func (s *Service) CreateOwnerAccount(ctx context.Context) error {

	// Hash the password before proceeding
	modified, err := s.hasher.HashPlain(s.password)
	if err != nil {
		return err
	}

	u := user.User{
		ID:       primitive.NewObjectID(),
		FullName: "Aymen DHAHRI",
		Email:    s.contactEmail,
		Password: modified,
		Photo:    "https://github.com/DroidZed.png",
	}

	s.logger.Debug("body", slog.Any("here", u))

	// GetOne when no doc fount: return back an err indicating no doc found
	found, _ := s.userSrv.GetOne(
		ctx,
		utils.GetInput{
			Filter: bson.M{"email": u.Email},
		},
	)

	// if the doc has been found, ignore and return
	if found != nil {
		return nil
	}

	// Else we insert the user on an empty collection
	if err := s.userSrv.Add(ctx, u); err != nil {
		return err
	}

	s.logger.Info("admin created.")

	return nil
}
