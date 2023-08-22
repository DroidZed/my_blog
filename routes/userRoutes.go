package routes

import (
	"github.com/DroidZed/go_lance/controllers"
	"github.com/go-chi/chi"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Get("/", controllers.GetAllUsers)

	userRouter.Route("/{id}", func(r chi.Router) {
		r.Get("/", controllers.GetUserById)
		r.Delete("/", controllers.DeleteUserById)
	})

	return userRouter
}
