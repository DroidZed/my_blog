package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/DroidZed/go_lance/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	loginBody := &LoginBody{}

	if err := json.NewDecoder(r.Body).Decode(loginBody); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, LoginResponse{Error: err.Error()})
		return
	}

	userId, err := ValidateUser(loginBody)

	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: err.Error()})
		return
	}

	access, err := cryptor.GenerateAccessToken(userId)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: err.Error()})
		return
	}

	refresh, err := cryptor.GenerateRefreshToken()
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: err.Error()})
		return
	}

	loginResp := LoginResponse{
		Jwt:     access,
		Refresh: refresh,
	}

	utils.JsonResponse(w, http.StatusOK, loginResp)

}

// TODO: Complete this function!
func Register(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	user := &user.User{ID: primitive.NewObjectID()}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, LoginResponse{Error: err.Error()})
		return
	}

	userId, err := SaveUser(user)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusCreated, userId)

}
