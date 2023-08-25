package controllers

// CRUD: user

import (
	"encoding/json"
	"fmt"
	"github.com/DroidZed/go_lance/config"
	"github.com/DroidZed/go_lance/models"
	"github.com/DroidZed/go_lance/services"
	"github.com/DroidZed/go_lance/utils"
	"github.com/ggicci/httpin"
	"net/http"
)

type UserController interface {
	GetAllUsers(w http.ResponseWriter, _ *http.Request)
	GetUserById(w http.ResponseWriter, _ *http.Request)
	DeleteUserById(w http.ResponseWriter, r *http.Request)
	UpdateUserById(w http.ResponseWriter, r *http.Request)
}

func GetAllUsers(w http.ResponseWriter, _ *http.Request) {

	results := services.FindAllUsers()

	if results == nil {
		utils.JsonResponse(w, http.StatusInternalServerError, []models.User{})
		return
	}

	utils.JsonResponse(w, http.StatusOK, results)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(httpin.Input).(*UserIdPath).UserId

	result := services.FindUserById(id)

	if result == nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, result)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(httpin.Input).(*UserIdPath).UserId

	result := services.DeleteOne(id)

	if !result {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: fmt.Sprintf("User with id: %s has been deleted successfully!", id)})
}

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(httpin.Input).(*UserIdPath).UserId

	log := config.Logger.LogHandler

	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Fatal(err)
	}

	err := services.UpdateOne(id, user)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid update!"})
	}
}
