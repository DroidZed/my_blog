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
	"github.com/DroidZed/my_blog/internal/views"
	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	Router    *chi.Mux
	EnvConfig *config.EnvConfig
}

type ServerDefinition interface {
	New() (server *Server)
	ApplyMiddleWares()
	MountHandlers()
	InitOwner()
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
		views.NotFound().Render(r.Context(), w)
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

func (*Server) InitOwner() {

	log := config.GetLogger()

	userService := &user.UserService{}

	user := &user.User{
		ID:       primitive.NewObjectID(),
		FullName: "Aymen DHAHRI",
		Email:    config.LoadEnv().MASTER_EMAIL,
		Password: config.LoadEnv().MASTER_PWD,
		Photo:    "https://github.com/DroidZed.png",
	}

	if found := userService.FindUserByEmail(user.Email); found != nil {
		return
	}

	if err := userService.SaveUser(user); err != nil {
		log.Fatalf("Error occurred while saving the user to the db\n %s", err.Error())
	}

	log.Info("Admin created with password from env.")
}
