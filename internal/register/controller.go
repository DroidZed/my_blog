package register

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/DroidZed/go_lance/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Complete this function!
func Register(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	user := &user.User{ID: primitive.NewObjectID()}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	userId, err := SaveUser(user)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusCreated, userId)

}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {

}
