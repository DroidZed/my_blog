package user

import (
	md "github.com/DroidZed/go_lance/internal/middleware"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Use(middleware.AllowContentType("application/json"))

	userRouter.With(md.AccessVerify).Group(func(r chi.Router) {
		r.Get("/", GetAllUsers)
		r.Put("/", UpdateUser)
		r.With(httpin.NewInput(UserIdPath{})).Group(func(ru chi.Router) {
			ru.Get("/{userId}", GetUserById)
			ru.Delete("/{userId}", DeleteUserById)
		})
	})
	return userRouter
}
