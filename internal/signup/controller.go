package signup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/pigeon"
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
	_, err := signupService.SaveUser(u)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s", err.Error())})
		return
	}

	confirmationCodeE := &ConfirmationCode{
		ID:        primitive.NewObjectID(),
		Code:      generateCode(999999),
		Email:     u.Email,
		ExpiresAt: primitive.NewDateTimeFromTime(time.Now().Add(time.Duration(time.Duration.Minutes(15)))),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	result, err := signupService.SaveConfirmationCode(confirmationCodeE)
	if err != nil {
		user.DeleteUserById(w, r)
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Error creating the user!\n%s", err.Error())})
		return
	}

	deliverEmailToUser(w, r, u.Email, "CONFIRMATION MAIL", "confirmation_email", "Code has been delivered to your inbox.", result)
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {

}

func deliverEmailToUser(w http.ResponseWriter, r *http.Request,
	to,
	subject,
	templateName,
	response string,
	data interface{},
) {
	log := config.InitializeLogger().LogHandler

	req := pigeon.NewRequest(
		[]string{to},
		subject,
		"",
	)

	err := req.ParseTemplate(templateName, data)

	if err != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	emailErr := req.SendEmail()

	if emailErr != nil {
		log.Error(err)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: err.Error()})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: response})
}

func generateCode(bound int64) string {
	builder := &strings.Builder{}

	for i := 0; i < 6; i++ {
		builder.WriteString(fmt.Sprintf("%d", utils.RNG(bound)))
	}

	return builder.String()
}
