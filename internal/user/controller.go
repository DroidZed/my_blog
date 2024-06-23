package user

// CRUD: user

import (
	"fmt"
	"net/http"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/middleware"
	"github.com/DroidZed/my_blog/internal/utils"
)

func GetUserById(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value(middleware.AuthCtxKey{}).(string)

	log := config.GetLogger()

	userService := &UserService{}

	user, err := userService.FindUserByID(id)

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	log := config.GetLogger()

	var user User

	if err := utils.DecodeBody(r, &user); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	log.Debug(user)

	userService := &UserService{}

	if err := userService.SaveUser(&user); err != nil {
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Invalid update!\n %s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: "Success !!"})
}
