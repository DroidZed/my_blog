package user

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type UserProvider interface {
	Add(ctx context.Context, u User) error
	GetByIDProj(ctx context.Context, id string, in utils.GetInput) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetOne(ctx context.Context, in utils.GetInput) (*User, error)
	Validate(ctx context.Context, email, password string) (*User, error)
}

type Controller struct {
	service UserProvider
	logger  *slog.Logger
}

func NewController(up UserProvider, l *slog.Logger) Controller {
	return Controller{
		service: up,
		logger:  l,
	}
}

func (c Controller) GetUserByID(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(jwtverify.AuthCtxKey{}).(string)

	user, err := c.service.GetByIDProj(
		r.Context(),
		id,
		utils.GetInput{
			Projection: bson.D{
				{Key: "password", Value: 0},
			},
		},
	)

	if err != nil {
		c.logger.Error("failed to find user", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}
