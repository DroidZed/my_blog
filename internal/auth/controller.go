package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/utils"
)

type AuthService interface {
	GenerateNewTokens(expiredToken string) (string, string, error)
	CreateLoginResponse(ctx context.Context, body LoginBody) (LoginResponse, error)
	CreateOwnerAccount(ctx context.Context) error
}

type Controller struct {
	authSrv AuthService
	logger  *slog.Logger
	hasher  cryptor.CryptoHelper
}

func NewController(
	authSrv AuthService,
	logger *slog.Logger,
	hasher cryptor.CryptoHelper,
) *Controller {
	return &Controller{
		authSrv: authSrv,
		logger:  logger,
		hasher:  hasher,
	}
}

// LoginReq Auth
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
func (c *Controller) LoginReq(w http.ResponseWriter, r *http.Request) {

	var loginBody LoginBody

	if err := utils.DecodeBody(r, &loginBody); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	resp, err := c.authSrv.CreateLoginResponse(r.Context(), loginBody)

	if err != nil {
		c.logger.Error("login", slog.String("err", err.Error()))
	}

	utils.JsonResponse(w, http.StatusOK, resp)
}

// RefreshTheAccessToken Refresh
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
func (c *Controller) RefreshTheAccessToken(w http.ResponseWriter, r *http.Request) {

	var refreshReq RefreshReq

	if err := utils.DecodeBody(r, &refreshReq); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	newAcc, newRef, err := c.authSrv.GenerateNewTokens(refreshReq.Expired)

	if err != nil {
		c.logger.Error("generating new tokens", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: "Unable to process your request"})
		return
	}

	utils.JsonResponse(w, http.StatusOK, LoginResponse{Jwt: newAcc, Refresh: newRef})
}
