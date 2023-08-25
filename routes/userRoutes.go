package routes

import (
	"github.com/DroidZed/go_lance/controllers"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Get("/", controllers.GetAllUsers)

	userRouter.Route("/{id}", func(r chi.Router) {
		r.Get("/", controllers.GetUserById)
		r.Put("/", controllers.UpdateUserById)
		r.Delete("/", controllers.DeleteUserById)
	})

	// 	r.With(httpin.NewInput(controllers.UserIdParam{}))

	return userRouter
}
