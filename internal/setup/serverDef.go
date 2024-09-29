package setup

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/DroidZed/my_blog/cmd/web"
	"github.com/DroidZed/my_blog/cmd/web/pages"
	_ "github.com/DroidZed/my_blog/docs"
	"github.com/a-h/templ"
	_ "github.com/joho/godotenv/autoload"

	"github.com/DroidZed/my_blog/internal/database"
	"github.com/DroidZed/my_blog/internal/httpslog"
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
	GetArticleInfo(w http.ResponseWriter, r *http.Request)
	AddArticle(w http.ResponseWriter, r *http.Request)
}

type UserManager interface {
	GetUserByID(w http.ResponseWriter, r *http.Request)
}

type AuthManager interface {
	LoginReq(w http.ResponseWriter, r *http.Request)
	RefreshTheAccessToken(w http.ResponseWriter, r *http.Request)
}

type JwtMiddleware interface {
	AccessVerify(next http.Handler) http.Handler
	RefreshVerify(next http.Handler) http.Handler
}

type Server struct {
	logger *slog.Logger

	articleManager ArticleManager
	userManager    UserManager
	authManager    AuthManager

	authMiddleware JwtMiddleware

	database *database.Service
}

var (
	port string = os.Getenv("PORT")
	host string = os.Getenv("HOST")
)

func NewServer(
	logger *slog.Logger,
	articleManager ArticleManager,
	userManager UserManager,
	authManager AuthManager,
	authMiddleware JwtMiddleware,
	database *database.Service,
) *Server {
	return &Server{
		logger:         logger,
		articleManager: articleManager,
		authManager:    authManager,
		userManager:    userManager,
		authMiddleware: authMiddleware,
		database:       database,
	}
}

func (s *Server) MountHandlers(r *chi.Mux) {

	r.NotFound(templ.Handler(pages.NotFound()).ServeHTTP)

	r.Get("/", templ.Handler(pages.Index()).ServeHTTP)

	r.Get("/login", templ.Handler(pages.Login()).ServeHTTP)

	r.Handle("/assets/*", http.FileServer(http.FS(web.Assets)))

	r.Get("/api/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("http://%s/api/swagger/doc.json",
				net.JoinHostPort(
					host,
					port,
				),
			)),
	))

	// Article
	r.Route("/articles", func(r chi.Router) {
		r.Get("/{title}", s.articleManager.GetArticle)
		r.Get("/info/{id}", s.articleManager.GetArticleInfo)
	})

	// Auth
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", s.authManager.LoginReq)
		r.With(s.authMiddleware.RefreshVerify).Group(func(r chi.Router) {
			r.Post("/refresh-token", s.authManager.RefreshTheAccessToken)
		})
	})

	// User
	r.Route("/api/user", func(r chi.Router) {
		r.Use(s.authMiddleware.AccessVerify)
		r.Get("/", s.userManager.GetUserByID)
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
