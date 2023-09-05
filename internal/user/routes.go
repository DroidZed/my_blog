package user

import (
	customMiddleware "github.com/DroidZed/go_lance/internal/middleware"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Use(middleware.AllowContentType("application/json"))
	userRouter.Use(customMiddleware.JwtVerify)

	userRouter.Route("/", func(r chi.Router) {
		r.Get("/", GetAllUsers)
		r.Put("/", UpdateUser)

	})

	userRouter.With(httpin.NewInput(utils.UserIdPath{})).
		Route("/{userId}", func(r chi.Router) {
			r.Get("/", GetUserById)
			r.Delete("/", DeleteUserById)
		})

	return userRouter
}
