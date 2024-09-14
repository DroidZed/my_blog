package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)



type Controller struct {
	service    UserProvider
	logger *slog.Logger
}

func NewController(up UserProvider, l *slog.Logger) Controller {
	return Controller{
		service:    up,
		logger: l,
	}
}

func (c Controller) GetUserById(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(jwtverify.AuthCtxKey{}).(string)

	user, err := c.service.GetByIdProj(
		r.Context(),
		id,
		bson.D{{Key: "password", Value: 0}},
	)

	if err != nil {
		c.logger.Error("failed to find user", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}
