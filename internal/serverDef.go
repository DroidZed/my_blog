package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"path/filepath"

	_ "github.com/DroidZed/go_lance/docs"
	"github.com/DroidZed/go_lance/internal/auth"
	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/pigeon"
	"github.com/DroidZed/go_lance/internal/signup"
	"github.com/DroidZed/go_lance/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/withmandala/go-log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Router    *chi.Mux
	DbClient  *mongo.Client
	EnvConfig *config.EnvConfig
	Logger    *log.Logger
}

func CreateNewServer() *Server {
	server := &Server{}
	server.Logger = config.InitializeLogger().LogHandler
	server.Router = chi.NewRouter()
	server.DbClient = config.GetConnection()
	server.EnvConfig = config.LoadEnv()
	pigeon.GetSmtp()
	return server
}

func (s *Server) MountViewsFolder() {
	if workDir, err := os.Getwd(); err != nil {
		s.Logger.Errorf(err.Error())
	} else {
		filesDir := http.Dir(filepath.Join(workDir, "public/views"))
		FileServer(s.Router, "/", filesDir)
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (s *Server) MountHandlers() {

	// Mount all handlers here
	s.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.EnvConfig.Port)), //The url pointing to API definition
	))

	s.Router.Mount("/signup", signup.SignUpRoutes())
	s.Router.Mount("/auth", auth.AuthRoutes())
	s.Router.Mount("/user", user.UserRoutes())
}

func (s *Server) ApplyMiddleWares() {

	s.Router.Use(middleware.RequestID)

	s.Router.Use(middleware.CleanPath)

	s.Router.Use(middleware.URLFormat)

	s.Router.Use(middleware.StripSlashes)

	s.Router.Use(httplog.RequestLogger(httplog.NewLogger("GoLance-Log", httplog.Options{
		// JSON:             false,
		LogLevel:        slog.LevelDebug,
		Concise:         true,
		RequestHeaders:  false,
		TimeFieldFormat: time.DateTime,
		Tags: map[string]string{
			"env": s.EnvConfig.Env,
		},
	})))

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
