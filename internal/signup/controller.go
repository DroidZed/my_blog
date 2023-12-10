package signup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/user"

	"github.com/DroidZed/go_lance/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Complete this function!
func Register(w http.ResponseWriter, r *http.Request) {

	log := config.InitializeLogger().LogHandler

	u := &user.User{ID: primitive.NewObjectID()}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	signupService := &SignUpService{}

	userId, err := signupService.SaveUser(u)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s", err.Error())})
		return
	}

	builder := &strings.Builder{}

	for i := 0; i < 6; i++ {

		builder.WriteString(fmt.Sprintf("%d", utils.RNG(999999)))
	}

	confirmationCodeE := &ConfirmationCode{
		ID:    primitive.NewObjectID(),
		Code:  builder.String(),
		Email: u.Email,
	}

	result, err := signupService.SaveConfirmationCode(confirmationCodeE)
	if err != nil {
		user.DeleteUserById(w, r)
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s", err.Error())})
		return
	}

	log.Debug(result)

	utils.JsonResponse(w, http.StatusCreated, userId)

}

// TODO: Figure it out :p
func VerifyEmail(w http.ResponseWriter, r *http.Request) {

}

func TestSendingAnEmail(w http.ResponseWriter, r *http.Request) {

}
