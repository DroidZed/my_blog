package routes

import (
	controller "github.com/DroidZed/go_lance/src/controllers"
	"github.com/go-chi/chi"

)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Get("/", controller.GetAllUsers)

	userRouter.Route("/{id}", func(r chi.Router) {
		r.Get("/", controller.GetUserById)
		r.Delete("/", controller.DeleteUserById)
	})

	return userRouter
}
