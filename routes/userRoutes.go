package routes

import (
	controller "github.com/DroidZed/go_lance/controllers"
	"github.com/go-chi/chi"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Get("/", controller.GetAllUsers)
	userRouter.Get("/{id}", controller.GetUserById)

	return userRouter
}
