package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DroidZed/my_blog/internal/auth"
	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/forgotPwd"
	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/pigeon"
	"github.com/DroidZed/my_blog/internal/setup"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/go-chi/chi/v5"
)

func startService(
	ctx context.Context,
	mux *chi.Mux,
	logger *slog.Logger,
) (*setup.Server, error) {

	env, err := config.LoadEnv()
	if err != nil {
		return nil, err
	}

	pigeon := pigeon.Pigeon{
		Auth: smtp.PlainAuth(
			"pigeon",
			env.SmtpUsername,
			env.SmtpPassword,
			env.SmtpHost,
		),
		From: env.SmtpUsername,
		Addr: net.JoinHostPort(env.SmtpHost, env.SmtpPort),
	}

	cHelper := &cryptor.Cryptor{
		AccessExpiry:  env.AccessExpiry,
		AccessSecret:  env.AccessSecret,
		RefreshExpiry: env.RefreshExpiry,
		RefreshSecret: env.RefreshSecret,
	}

	logger.Info("opening a database connection...")

	dbClient, err := config.GetConnection(ctx, env)
	if err != nil {
		return nil, err
	}

	logger.Info("connected to ", slog.String("dbName", env.DBName))

	db := dbClient.Database(env.DBName)

	userService := &user.Service{
		Hasher: cHelper,
		Db:     db,
	}

	authController := &auth.Controller{
		UserService:   userService,
		Logger:        logger,
		CHelper:       cHelper,
		MASTER_EMAIL:  env.MASTER_EMAIL,
		MASTER_PWD:    env.MASTER_PWD,
		RefreshSecret: env.RefreshSecret,
	}

	server := &setup.Server{
		Authenticator: authController,
		UserManager: &user.Controller{
			UserService: userService,
			Logger:      logger,
		},
		PasswordManager: &forgotPwd.Controller{
			UserService: userService,
			Pigeon:      pigeon,
			Logger:      logger,
		},
		AuthMiddleware: jwtverify.JwtVerify{
			Logger:        logger,
			CHelper:       cHelper,
			AccessSecret:  env.AccessSecret,
			RefreshSecret: env.RefreshSecret,
		},
	}

	envPort := server.EnvConfig.Port

	server.ApplyMiddleWares(mux)

	server.MountHandlers(mux)

	authController.InitOwner(ctx)

	logger.Info("listening", slog.Int64("port", envPort))

	return server, nil
}

// Entry point, setting up chi and graceful shutdown <3
// @title My Website's API Docs
// @version 1.0
// @description This is the GoLance API documentation.
// @termsOfService https://droidzed.tn/terms/

// @contact.name Aymen Dhahri
// @contact.url https://droidzed.tn/support
// @contact.email droid.zed77@outlook.com

// @license.name MIT
// @license.url https://github.com/DroidZed/my_blog/LICENSE

// @host droidzed.tn
// @BasePath /
func main() {
	code := run()

	os.Exit(code)
}

func run() int {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	mux := chi.NewRouter()

	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{})

	l := slog.New(logHandler)

	app, err := startService(ctx, mux, l)
	if err != nil {
		l.Error("error with server startup", slog.String("err", err.Error()))
		return -1
	}

	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.EnvConfig.Port),
		Handler:  mux,
		ErrorLog: slog.NewLogLogger(l.Handler(), slog.LevelError),
	}

	// Run the server in a separate goroutine.
	go func() {
		defer cancel()

		// Ignore the error returned because it'll always be [http.ErrServerClosed].
		_ = server.ListenAndServe()
	}()

	// Wait for context to be cancelled.
	<-ctx.Done()

	// Start the shutdown procedure.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Error("something went wrong during shutdown",
			slog.Any("error", err),
		)

		return -1
	}

	l.Info("goodbye 🧩 👋")
	return 0
}
