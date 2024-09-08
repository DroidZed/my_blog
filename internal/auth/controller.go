package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	UserService user.UserProvider
	Logger      *slog.Logger
	CHelper     cryptor.CryptoHelper

	MASTER_EMAIL  string
	MASTER_PWD    string
	RefreshSecret string
}

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
func (c *Controller) LoginReq(w http.ResponseWriter, r *http.Request) {

	var loginBody LoginBody

	if err := utils.DecodeBody(r, &loginBody); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	user, err := c.validateUser(r.Context(), &loginBody)

	if err != nil {
		c.Logger.Error("validation failed with", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: "Invalid credentials"})
		return
	}

	userId := user.ID.Hex()

	access, err := c.CHelper.GenerateAccessToken(userId)
	if err != nil {
		c.Logger.Error("token gen failed with", slog.String("err", err.Error()))
	}

	refresh, err1 := c.CHelper.GenerateRefreshToken()
	if err1 != nil {
		c.Logger.Error("token gen failed with", slog.String("err", err1.Error()))
	}

	if access == "" || refresh == "" {
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, LoginResponse{
		Jwt:     access,
		Refresh: refresh,
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
func (c *Controller) RefreshTheAccessToken(w http.ResponseWriter, r *http.Request) {

	var refreshReq RefreshReq

	if err := utils.DecodeBody(r, &refreshReq); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid JSON payload"})
		return
	}

	expiredToken := refreshReq.Expired

	access, err := c.CHelper.ParseToken(expiredToken, c.RefreshSecret)
	if err != nil {
		c.Logger.Error("token gen failed with", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	userId, err := c.CHelper.ExtractSubFromClaims(access)
	if err != nil {
		c.Logger.Error("claims corrupted", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	newAcc, err := c.CHelper.GenerateAccessToken(userId)
	if err != nil {
		c.Logger.Error("token gen failed with", slog.String("err", err.Error()))

	}

	newRef, err1 := c.CHelper.GenerateRefreshToken()
	if err1 != nil {
		c.Logger.Error("token gen failed with", slog.String("err", err1.Error()))
	}

	if newAcc != "" || newRef != "" {
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, LoginResponse{Jwt: newAcc, Refresh: newRef})
}

func (c *Controller) InitOwner(ctx context.Context) error {

	u := &user.User{
		ID:       primitive.NewObjectID(),
		FullName: "Aymen DHAHRI",
		Email:    c.MASTER_EMAIL,
		Password: c.MASTER_PWD,
		Photo:    "https://github.com/DroidZed.png",
	}

	found, err := c.UserService.GetOne(ctx, bson.M{"email": u.Email}, nil)

	if err != nil {
		return err
	}

	if found != nil {
		return nil
	}

	if err := c.UserService.Add(ctx, u); err != nil {
		return err
	}

	c.Logger.Info("admin created.")
	return nil
}

func (c *Controller) validateUser(ctx context.Context, login *LoginBody) (*user.User, error) {

	data, err := c.UserService.GetOne(ctx, bson.M{"email": login.Email}, nil)
	if err != nil {
		return nil, err
	}

	pwdIsValid := c.CHelper.CompareSecureToPlain(data.Password, login.Password)

	if !pwdIsValid {
		return nil, fmt.Errorf("invalid credentials")
	}

	return data, nil
}
