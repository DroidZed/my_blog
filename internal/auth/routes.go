package auth

import (
	md "github.com/DroidZed/go_lance/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {

	authRouter := chi.NewRouter()

	authRouter.Route("/", func(r chi.Router) {

		r.Post("/login", Login)
	})

	authRouter.Route("/refresh-token", func(r chi.Router) {
		r.Use(md.RefreshVerify)
		r.Post("/", RefreshTheAccessToken)
	})

	return authRouter

}
