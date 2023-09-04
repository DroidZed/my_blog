package auth

import (
	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {

	authRouter := chi.NewRouter()

	authRouter.Route("/", func(r chi.Router) {

		r.Post("/login", Login)
		r.Post("/register", Register)
	})

	return authRouter

}
