package user

// CRUD: user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/utils"
)

type Controller struct {
	UserService *Service
	Logger      *slog.Logger
}

func (c *Controller) GetUserById(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(jwtverify.AuthCtxKey{}).(string)

	log := config.GetLogger()

	user, err := c.UserService.FindUserByID(id)

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {

	log := config.GetLogger()

	var user User

	if err := utils.DecodeBody(r, &user); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	log.Debug(user)

	if err := c.UserService.SaveUser(&user); err != nil {
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Invalid update!\n %s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: "Success !!"})
}
