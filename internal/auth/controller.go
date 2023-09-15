package auth

import (
	"encoding/json"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/cryptor"
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
	conf := config.LoadConfig()

	refreshBody := &RefreshReq{}

	if err := json.NewDecoder(r.Body).Decode(refreshBody); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	expiredToken := refreshBody.Expired

	access, err := cryptor.ParseToken(expiredToken, conf.AccessSecret)
	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	userId, err := cryptor.ExtractSubFromClaims(access)
	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	newAcc, err := cryptor.GenerateAccessToken(userId)
	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	newRef, err := cryptor.GenerateRefreshToken()
	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, JwtResponse{Jwt: newAcc, Refresh: newRef})
}
