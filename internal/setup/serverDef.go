package setup

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"

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

type ArticleManager interface {
	GetArticle(w http.ResponseWriter, r *http.Request)
	AddArticle(w http.ResponseWriter, r *http.Request)
}

type Authenticator interface {
	LoginReq(w http.ResponseWriter, r *http.Request)
	RefreshTheAccessToken(w http.ResponseWriter, r *http.Request)
}

type UserManager interface {
	GetUserByID(w http.ResponseWriter, r *http.Request)
}

type JwtMiddleware interface {
	AccessVerify(next http.Handler) http.Handler
	RefreshVerify(next http.Handler) http.Handler
}

type Server struct {
	env    *config.EnvConfig
	logger *slog.Logger

	authProvider   Authenticator
	userProvider   UserManager
	articleManager ArticleManager

	authMiddleware JwtMiddleware
}

func NewServer(
	env *config.EnvConfig,
	logger *slog.Logger,
	authMiddleware JwtMiddleware,
	auth Authenticator,
	userProvider UserManager,
	articleManager ArticleManager,
) *Server {
	return &Server{
		env:            env,
		authProvider:   auth,
		logger:         logger,
		userProvider:   userProvider,
		authMiddleware: authMiddleware,
		articleManager: articleManager,
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
			fmt.Sprintf("http://%s/api/swagger/doc.json",
				net.JoinHostPort(
					s.env.Host,
					strconv.FormatInt(s.env.Port, 10),
				),
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
		r.Get("/", s.userProvider.GetUserByID)
	})

	// Article
	r.Route("/articles", func(r chi.Router) {
		// r.With(s.authMiddleware.AccessVerify).Group(func(r chi.Router) {})
		r.Post("/", s.articleManager.AddArticle)
		r.Get("/{id}", s.articleManager.GetArticle)
	})
}

func (s *Server) ApplyMiddleWares(r *chi.Mux) {
	r.Use(middleware.RequestID)

	r.Use(middleware.URLFormat)

	r.Use(middleware.StripSlashes)

	r.Use(middleware.CleanPath)

	r.Use(httpslog.New(s.logger))

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
