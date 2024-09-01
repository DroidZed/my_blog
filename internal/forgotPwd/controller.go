package forgotPwd

import (
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/user"
	"github.com/DroidZed/my_blog/internal/utils"
)

type ForgotPwd struct {
	UserService *user.Service
	Logger      *slog.Logger
}

func (f ForgotPwd) DoSendMagicLink(w http.ResponseWriter, r *http.Request) {

	var body ForgotPwdReq

	err := utils.DecodeBody(r, &body)

	if err != nil {
		f.Logger.Error("invalid body", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: "Error parsing request body, aborting."})
		return
	}

	// resend := r.URL.Query().Get("resend")

	user := f.UserService.FindUserByEmail(body.Email)

	if user == nil {
		f.Logger.Error("invalid user", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid email"})
		return
	}
}

func (f ForgotPwd) DoValidateMagicLink(w http.ResponseWriter, r *http.Request) {

}

func generateAndSendMagicLink(email string) error {
	return nil
}
