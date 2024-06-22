package auth

import (
	"net/http"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/utils"
)

// LoginReq Auth godoc
//
//	@Summary		Auth user
//	@Description	Get token, user basic data
//	@Tags			auth
//	@Accept			json
//	@Param			login	body		LoginBody	true	"Login User"
//	@Produce		json
//	@Success		200	{object}	LoginResponse
//	@Failure		404	{object}	LoginResponse
//	@Router			/api/auth/login [post]
func LoginReq(w http.ResponseWriter, r *http.Request) {

	log := config.GetLogger()

	var loginBody LoginBody

	if err := utils.DecodeBody(r, &loginBody); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	user, err := ValidateUser(&loginBody)

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: "Invalid credentials"})
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

// RefreshTheAccessToken Refresh godoc
//
//	@Summary		Refresh tokens
//	@Description	Refresh access + refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LoginBody
//	@Failure		404	{object}	LoginResponse
//	@Security		Bearer
//	@Router			/api/auth/refresh-token [post]
func RefreshTheAccessToken(w http.ResponseWriter, r *http.Request) {

	log := config.GetLogger()

	conf := config.LoadEnv()

	var refreshReq RefreshReq

	if err := utils.DecodeBody(r, &refreshReq); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	expiredToken := refreshReq.Expired

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
