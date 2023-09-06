package auth

import (
	"github.com/DroidZed/go_lance/internal/auth/login"
	"github.com/DroidZed/go_lance/internal/auth/register"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {

	authRouter := chi.NewRouter()

	authRouter.Route("/", func(r chi.Router) {

		r.Post("/login", login.Login)
		r.Post("/register", register.Register)
	})

	return authRouter

}
