package forgotPwd

import (
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/pigeon"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/DroidZed/my_blog/internal/utils"
)

type Controller struct {
	UserService user.UserProvider
	Pigeon      *pigeon.Pigeon
	Logger      *slog.Logger
}

func (c Controller) DoSendMagicLink(w http.ResponseWriter, r *http.Request) {

	var body ForgotPwdReq

	err := utils.DecodeBody(r, &body)

	if err != nil {
		c.Logger.Error("invalid body", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: "Error parsing request body, aborting."})
		return
	}

	// resend := r.URL.Query().Get("resend")

	user := c.UserService.FindUserByEmail(body.Email)

	if user == nil {
		c.Logger.Error("invalid user", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: "Invalid email"})
		return
	}
}

func (c Controller) DoValidateMagicLink(w http.ResponseWriter, r *http.Request) {

}

func (c Controller) generateAndSendMagicLink(email string) error {
	return nil
}
