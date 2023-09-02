package user

// CRUD: user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/ggicci/httpin"
)

type UserController interface {
	GetAllUsers(w http.ResponseWriter, _ *http.Request)
	GetUserById(w http.ResponseWriter, _ *http.Request)
	DeleteUserById(w http.ResponseWriter, r *http.Request)
	UpdateUserById(w http.ResponseWriter, r *http.Request)
}

func GetAllUsers(w http.ResponseWriter, _ *http.Request) {

	results := FindAllUsers()

	if results == nil {
		utils.JsonResponse(w, http.StatusInternalServerError, []User{})
		return
	}

	utils.JsonResponse(w, http.StatusOK, results)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(httpin.Input).(*utils.UserIdPath).UserId

	result := FindUserById(id)

	if result == nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, result)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(httpin.Input).(*utils.UserIdPath).UserId

	result := DeleteOne(id)

	if !result {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: fmt.Sprintf("User with id: %s has been deleted successfully!", id)})
}

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(httpin.Input).(*utils.UserIdPath).UserId

	log := config.InitializeLogger().LogHandler

	user := &User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Fatal(err)
	}

	err := UpdateOne(id, user)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid update!"})
	}
}
