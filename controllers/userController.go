package controllers

// CRUD: user

import (
	"fmt"
	"github.com/DroidZed/go_lance/models"
	"github.com/DroidZed/go_lance/services"
	"github.com/DroidZed/go_lance/utils"
	"github.com/go-chi/chi"
	"net/http"
)

func GetAllUsers(w http.ResponseWriter, _ *http.Request) {

	results := services.FindAll()

	if results == nil {
		utils.JsonResponse(w, http.StatusInternalServerError, []models.User{})
		return
	}

	utils.JsonResponse(w, http.StatusOK, results)

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result := services.FindOneById(id)

	if result == nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, result)

}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result := services.DeleteOne(id)

	if !result {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("User with id %s could not be found.", id)})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: fmt.Sprintf("User with id: %s has been deleted successfully!", id)})
}
