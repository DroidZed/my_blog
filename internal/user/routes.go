package user

import (
	customMiddleware "github.com/DroidZed/go_lance/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Use(middleware.AllowContentType("application/json"))

	userRouter.Group(func(r chi.Router) {
		r.Use(customMiddleware.AccessVerify)
		r.Get("/", GetAllUsers)
		r.Put("/", UpdateUser)
		r.Route("/{userId}", func(r chi.Router) {
			r.Get("/", GetUserById)
			r.Delete("/", DeleteUserById)
		})
	})

	return userRouter
}
