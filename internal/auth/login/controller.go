package login

import (
	"encoding/json"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
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

	access, refresh, err := GenerateLoginTokens(userId)

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, LoginResponse{Jwt: access, Refresh: refresh})

}
