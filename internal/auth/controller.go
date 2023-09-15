package auth

import (
	"encoding/json"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
)

// Auth godoc
//
//	@Summary		Auth user
//	@Description	Get token, user basic data
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LoginBody
//	@Failure		404	{object}	LoginResponse
//	@Router			/auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	loginBody := &LoginBody{}

	if err := json.NewDecoder(r.Body).Decode(loginBody); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, LoginResponse{Error: err.Error()})
		return
	}

	user, err := ValidateUser(loginBody)

	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: err.Error()})
		return
	}

	userId := user.ID.Hex()

	access, refresh, err := GenerateLoginTokens(userId)

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, LoginResponse{
		Jwt:       access,
		Refresh:   refresh,
		UserId:    userId,
		Role:      user.Role,
		AccStatus: user.AccStatus,
	})
}

// Refresh godoc
//
//	@Summary		Refresh tokens
//	@Description	Refresh access + refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LoginBody
//	@Failure		404	{object}	LoginResponse
//	@Security		Bearer
//	@Router			/auth/refresh-token [post]
func RefreshTheAccessToken(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	loginBody := &RefreshReq{}

	if err := json.NewDecoder(r.Body).Decode(loginBody); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}
}
