package setup

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	_ "github.com/DroidZed/my_blog/docs"
	"github.com/DroidZed/my_blog/internal/asset"
	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/httpslog"
	"github.com/DroidZed/my_blog/internal/views/pages"
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

type Server struct {
	env    *config.EnvConfig
	logger *slog.Logger

	articleManager ArticleManager
}

func NewServer(
	env *config.EnvConfig,
	logger *slog.Logger,
	articleManager ArticleManager,
) *Server {
	return &Server{
		env:            env,
		logger:         logger,
		articleManager: articleManager,
	}
}

func (s *Server) MountHandlers(r *chi.Mux) {

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		pages.NotFound().Render(r.Context(), w)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Index().Render(r.Context(), w)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		pages.Login().Render(r.Context(), w)
	})

	asset.Mount(r)

	r.Get("/api/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("http://%s/api/swagger/doc.json",
				net.JoinHostPort(
					s.env.Host,
					strconv.FormatInt(s.env.Port, 10),
				),
			)),
	))

	// Article
	r.Route("/articles", func(r chi.Router) {
		r.Get("/{title}", s.articleManager.GetArticle)
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
