package user

import (
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Get("/", GetAllUsers)

	userRouter.With(httpin.NewInput(utils.UserIdParam{})).
		Route("/{id}", func(r chi.Router) {
			r.Get("/", GetUserById)
			r.Put("/", UpdateUserById)
			r.Delete("/", DeleteUserById)
		})

	return userRouter
}
