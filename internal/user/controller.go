package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/utils"
)

type Controller struct {
	UserService UserProvider
	Logger      *slog.Logger
}

func (c *Controller) GetUserById(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(jwtverify.AuthCtxKey{}).(string)

	user, err := c.UserService.FindUserByID(r.Context(), id)

	if err != nil {
		c.Logger.Error("failed to find user", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user User

	if err := utils.DecodeBody(r, &user); err != nil {
		c.Logger.Error("failed decode body", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	if err := c.UserService.SaveUser(r.Context(), &user); err != nil {
		c.Logger.Error("failed to save user", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Invalid update!\n %s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: "Success !!"})
}
