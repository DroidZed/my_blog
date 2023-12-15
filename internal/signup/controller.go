package signup

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/user"

	"github.com/DroidZed/go_lance/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler
	signupService := &SignUpService{}
	userService := &user.UserService{}

	decodedUser, err := decodeBodyToUser(r)

	if err != nil {
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
	log.Debug(decodedUser.ID.String())

	found, _ := userService.FindUserByEmail(decodedUser.Email)

	if found != nil {
		utils.JsonResponse(
			w,
			http.StatusBadRequest,
			utils.DtoResponse{Error: "Already exists!"},
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

	if err := saveCodeAndSendEmail(signupService, decodedUser); err != nil {
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

}

func saveCodeAndSendEmail(
	signupService *SignUpService,
	u *user.User,
) error {

	code, err := signupService.SaveConfirmationCode(u.Email)

	if err != nil {
		return err
	}

	if emailErr := signupService.deliverEmailToUser(
		u.Email,
		"CONFIRMATION MAIL",
		"confirmation_email",
		code,
	); emailErr != nil {
		return emailErr
	}

	return nil
}

func decodeBodyToUser(r *http.Request) (*user.User, error) {
	u := &user.User{ID: primitive.NewObjectID()}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		return nil, err
	}

	return u, nil
}
