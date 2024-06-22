package internal

import (
	"fmt"
	"net/http"

	_ "github.com/DroidZed/my_blog/docs"
	"github.com/DroidZed/my_blog/internal/auth"
	"github.com/DroidZed/my_blog/internal/config"
	md "github.com/DroidZed/my_blog/internal/middleware"
	"github.com/DroidZed/my_blog/internal/pigeon"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	Router    *chi.Mux
	EnvConfig *config.EnvConfig
}

type ServerDefinition interface {
	New() (server *Server)
	ApplyMiddleWares()
	MountHandlers()
}

func (s *Server) New() (server *Server) {
	server = &Server{}
	server.Router = chi.NewRouter()
	server.EnvConfig = config.LoadEnv()
	pigeon.GetSmtp()
	return server
}

func (s *Server) MountHandlers() {
	s.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/views/404.tmpl")
	})

	s.Router.Get("/api/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("http://%s:%d/api/swagger/doc.json",
				s.EnvConfig.Host,
				s.EnvConfig.Port,
			)),
	))

	// Auth
	s.Router.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginReq)
		r.With(md.RefreshVerify).Group(func(r chi.Router) {
			r.Post("/refresh-token", auth.RefreshTheAccessToken)
		})
	})

	// User
	s.Router.Route("/api/user", func(r chi.Router) {
		r.Use(md.AccessVerify)
		r.Put("/", user.UpdateUser)
		r.Get("/", user.GetUserById)
	})
}

func (s *Server) ApplyMiddleWares() {
	s.Router.Use(middleware.StripSlashes)

	s.Router.Use(middleware.RequestID)

	s.Router.Use(middleware.CleanPath)

	s.Router.Use(middleware.URLFormat)

	s.Router.Use(httplog.LoggerWithName("GoLance-Log"))

	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.Router.Use(middleware.Heartbeat("/health"))

	s.Router.Use(middleware.Heartbeat("/ping"))
}
