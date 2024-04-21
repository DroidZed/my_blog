package signup

import (
	"fmt"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/user"

	"github.com/DroidZed/go_lance/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(w http.ResponseWriter, r *http.Request) {

	log := config.GetLogger()
	signupService := &SignUpService{}
	userService := &user.UserService{}

	decodedUser := &user.User{ID: primitive.NewObjectID()}

	if err := utils.DecodeBody(r, decodedUser); err != nil {
		log.Error(err)
		utils.JsonResponse(
			w,
			http.StatusBadRequest,
			utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s",
				err.Error(),
			)},
		)
		return
	}

	found := userService.FindUserByEmail(decodedUser.Email)

	if found != nil {
		utils.JsonResponse(
			w,
			http.StatusBadRequest,
			utils.DtoResponse{Error: "Email already used."},
		)
		return
	}

	if err := userService.SaveUser(decodedUser); err != nil {
		log.Error(err)
		utils.JsonResponse(
			w,
			http.StatusBadRequest,
			utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s",
				err.Error(),
			)},
		)
		return
	}

	if err := signupService.SaveCodeAndSendEmail(decodedUser.Email); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(
		w,
		http.StatusOK,
		utils.DtoResponse{Message: "Confirmation code has been delivered to your inbox."},
	)
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {

	log := config.GetLogger()
	signupService := &SignUpService{}
	userService := &user.UserService{}
	verifyCode := &VerifyCodeBody{}

	if err := utils.DecodeBody(r, verifyCode); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	u := userService.FindUserByEmail(verifyCode.Email)

	if u == nil {
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: "No user to be found."})
		return
	}

	if u.AccStatus == 1 {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Account already verified."})
		return
	}

	if !signupService.CheckCodeValidity(verifyCode.Email) {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Expired code, please resend for verification."})
		return
	}

	if err := userService.ActivateUserAccount(verifyCode.Email); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: "Account activated successfully, please login now!"})
}

func ResetVerifyCode(w http.ResponseWriter, r *http.Request) {

	resetBody := &ResetCodeBody{}
	signupService := &SignUpService{}
	log := config.GetLogger()

	if err := utils.DecodeBody(r, resetBody); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	if err := signupService.SaveCodeAndSendEmail(resetBody.Email); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Message: "Check your inbox for your new code!"})
}
