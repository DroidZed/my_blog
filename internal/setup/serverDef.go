package setup

import (
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/DroidZed/my_blog/docs"
	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/httpslog"
	"github.com/DroidZed/my_blog/internal/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type ServerDefinition interface {
	ApplyMiddleWares()
	MountHandlers()
}

type Authenticator interface {
	LoginReq(w http.ResponseWriter, r *http.Request)
	RefreshTheAccessToken(w http.ResponseWriter, r *http.Request)
}

type UserManager interface {
	GetUserById(w http.ResponseWriter, r *http.Request)
}

type JwtMiddleware interface {
	AccessVerify(next http.Handler) http.Handler
	RefreshVerify(next http.Handler) http.Handler
}

type Server struct {
	env *config.EnvConfig

	authProvider Authenticator
	userProvider UserManager

	authMiddleware JwtMiddleware
}

func NewServer(
	env *config.EnvConfig,
	auth Authenticator,
	userProvider UserManager,
	authMiddleware JwtMiddleware,
) *Server {
	return &Server{
		env:            env,
		authProvider:   auth,
		userProvider:   userProvider,
		authMiddleware: authMiddleware,
	}
}

func (s *Server) MountHandlers(r *chi.Mux) {

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		views.NotFound().Render(r.Context(), w)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		views.Login().Render(r.Context(), w)
	})

	r.Get("/api/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("http://%s:%d/api/swagger/doc.json",
				s.env.Host,
				s.env.Port,
			)),
	))

	// Auth
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", s.authProvider.LoginReq)
		r.With(s.authMiddleware.RefreshVerify).Group(func(r chi.Router) {
			r.Post("/refresh-token", s.authProvider.RefreshTheAccessToken)
		})
	})

	// User
	r.Route("/api/user", func(r chi.Router) {
		r.Use(s.authMiddleware.AccessVerify)
		r.Get("/", s.userProvider.GetUserById)
	})
}

func (s *Server) ApplyMiddleWares(r *chi.Mux) {
	r.Use(middleware.RequestID)

	r.Use(middleware.URLFormat)

	r.Use(middleware.StripSlashes)

	r.Use(middleware.CleanPath)

	// r.Use(httplog.LoggerWithName("GoLance-Log"))
	r.Use(httpslog.New(slog.Default()))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Heartbeat("/health"))

	r.Use(middleware.Heartbeat("/ping"))
}
