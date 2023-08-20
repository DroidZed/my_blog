package routes

import (
	controller "github.com/DroidZed/go_lance/controllers"
	"github.com/go-chi/chi"
)

func DummyRoute() chi.Router {
	userRouter := chi.NewRouter()

	userRouter.Get("/", controller.GetDummyEntity)

	return userRouter
}
