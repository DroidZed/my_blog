package tests

import (
	"github.com/DroidZed/go_lance/routes"
	"github.com/DroidZed/go_lance/utils"
	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type TestServer struct {
	Router *chi.Mux
}

func CreateNewTestServer() *TestServer {
	s := &TestServer{}
	s.Router = chi.NewRouter()
	return s
}

func (s *TestServer) MountTestHandlers() {
	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.CleanPath)
	s.Router.Use(middleware.URLFormat)
	s.Router.Use(httplog.LoggerWithName("CHI API"))

	// Mount all handlers here
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.JsonResponse(w, 200, utils.DtoResponse{Message: "Hello World!"})
	})

	s.Router.Get("/dev", func(w http.ResponseWriter, r *http.Request) {
		utils.LogAllRoutes(s.Router)
		utils.JsonResponse(w, 200, utils.DtoResponse{Message: "Nothing will be returned. This is just a dummy message. If you're a developer, check your console."})

	})

	s.Router.Mount("/user", routes.UserRoutes())
}
