package forgotPwd

import (
	"net/http"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/DroidZed/my_blog/internal/utils"
)

func DoSendMagicLink(w http.ResponseWriter, r *http.Request) {

	log := config.GetLogger()

	var body ForgotPwdReq

	err := utils.DecodeBody(r, &body)

	if err != nil {
		log.Errorf(body.Email)
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: "Error parsing request body, aborting."})
		return
	}

	// resend := r.URL.Query().Get("resend")

	userService := &user.UserService{}

	user := userService.FindUserByEmail(body.Email)

	if user == nil {
		log.Errorf("User not found for email: %s", body.Email)
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid email"})
		return
	}
}

func DoValidateMagicLink(w http.ResponseWriter, r *http.Request) {

}
