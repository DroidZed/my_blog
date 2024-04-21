package internal

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"path/filepath"

	_ "github.com/DroidZed/go_lance/docs"
	"github.com/DroidZed/go_lance/internal/auth"
	"github.com/DroidZed/go_lance/internal/config"
	md "github.com/DroidZed/go_lance/internal/middleware"
	"github.com/DroidZed/go_lance/internal/pigeon"
	"github.com/DroidZed/go_lance/internal/signup"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/MadAppGang/httplog"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/cors"
)

type Server struct {
	Router    *chi.Mux
	EnvConfig *config.EnvConfig
}

type ServerDefinition interface {
	New() (server *Server)
	MountViewsFolder() error
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

func (s *Server) MountViewsFolder() error {
	workDir, err := os.Getwd()

	if err != nil {
		return err
	}

	filesDir := http.Dir(filepath.Join(workDir, "public/views"))
	fileServerSetup(s.Router, "/", filesDir)

	return nil
}

func (s *Server) MountHandlers() {
	s.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("http://%s:%d/swagger/doc.json",
				s.EnvConfig.Host,
				s.EnvConfig.Port,
			)),
	))

	// Auth
	s.Router.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginReq)
		r.With(md.RefreshVerify).Group(func(r chi.Router) {
			r.Post("/refresh-token", auth.RefreshTheAccessToken)
		})
	})

	// Sign up
	s.Router.Route("/signup", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/", signup.Register)
		r.Post("/verify-email", signup.VerifyEmail)
		r.Put("/reset-code", signup.ResetVerifyCode)
	})

	// User
	s.Router.Route("/user", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(md.AccessVerify)
		r.Get("/", user.GetAllUsers)
		r.Put("/", user.UpdateUser)
		r.With(httpin.NewInput(user.UserIdPath{})).Group(func(ru chi.Router) {
			ru.Get("/{userId}", user.GetUserById)
			ru.Delete("/{userId}", user.DeleteUserById)
		})
	})
}

func (s *Server) ApplyMiddleWares() {

	s.Router.Use(md.FixUrl)

	s.Router.Use(middleware.RequestID)

	s.Router.Use(middleware.CleanPath)

	s.Router.Use(middleware.URLFormat)

	s.Router.Use(middleware.StripSlashes)

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

// fileServerSetup conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServerSetup(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		routeCtx := chi.RouteContext(r.Context())
		pathPrefix := string(strings.TrimSuffix(routeCtx.RoutePattern(), "/*"))
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
