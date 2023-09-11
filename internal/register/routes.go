package register

import (
	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {

	authRouter := chi.NewRouter()

	authRouter.Route("/", func(r chi.Router) {

		r.Post("/register", Register)
		r.Get("/verify-email", VerifyEmail)
	})

	return authRouter

}
