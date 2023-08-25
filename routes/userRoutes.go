package routes

import (
	"github.com/DroidZed/go_lance/controllers"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Get("/", controllers.GetAllUsers)

	userRouter.With(httpin.NewInput(controllers.UserIdParam{})).
		Route("/{id}", func(r chi.Router) {
			r.Get("/", controllers.GetUserById)
			r.Put("/", controllers.UpdateUserById)
			r.Delete("/", controllers.DeleteUserById)
		})

	return userRouter
}
