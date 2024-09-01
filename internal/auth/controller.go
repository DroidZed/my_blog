package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/DroidZed/my_blog/internal/utils"
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

	user, err := c.validateUser(&loginBody)

	if err != nil {
		c.Logger.Error("validation failed with", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, LoginResponse{Error: "Invalid credentials"})
		return
	}

	userId := user.ID.Hex()

	access, _ := c.CHelper.GenerateAccessToken(userId)
	refresh, _ := c.CHelper.GenerateRefreshToken()

	if access == "" || refresh == "" {
		c.Logger.Error("token gen failed with", slog.String("err", err.Error()))
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

	newAcc, _ := c.CHelper.GenerateAccessToken(userId)
	newRef, _ := c.CHelper.GenerateRefreshToken()

	if newAcc != "" || newRef != "" {
		c.Logger.Error("token gen failed with", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, LoginResponse{Jwt: newAcc, Refresh: newRef})
}

func (c *Controller) validateUser(login *LoginBody) (*user.User, error) {

	data := c.UserService.FindUserByEmail(login.Email)
	if data == nil {
		return nil, fmt.Errorf("no user found")
	}

	pwdIsValid := c.CHelper.CompareSecureToPlain(data.Password, login.Password)

	if !pwdIsValid {
		return nil, fmt.Errorf("invalid credentials")
	}

	return data, nil
}

func (c *Controller) InitOwner() error {

	user := &user.User{
		ID:       primitive.NewObjectID(),
		FullName: "Aymen DHAHRI",
		Email:    c.MASTER_EMAIL,
		Password: c.MASTER_PWD,
		Photo:    "https://github.com/DroidZed.png",
	}

	if found := c.UserService.FindUserByEmail(user.Email); found != nil {
		return fmt.Errorf("invalid email")
	}

	if err := c.UserService.SaveUser(user); err != nil {
		return err
	}

	c.Logger.Info("admin created.")
	return nil
}
