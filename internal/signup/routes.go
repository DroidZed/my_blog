package signup

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SignUpRoutes() chi.Router {

	signUpRoutes := chi.NewRouter()

	signUpRoutes.Use(middleware.AllowContentType("application/json"))

	signUpRoutes.Group(func(r chi.Router) {
		r.Post("/", Register)
		r.Post("/verify-email", VerifyEmail)
		r.Put("/reset-code", ResetVerifyCode)
	})

	return signUpRoutes
}
