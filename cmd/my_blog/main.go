package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DroidZed/my_blog/internal/article"
	"github.com/DroidZed/my_blog/internal/auth"
	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/setup"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func startService(
	ctx context.Context,
	mux *chi.Mux,
	env *config.EnvConfig,
	logger *slog.Logger,
) error {

	// pigeon := pigeon.New(
	// 	env.SmtpUsername,
	// 	env.SmtpPassword,
	// 	env.SmtpHost,
	// 	env.SmtpPort,
	// )

	cHelper := cryptor.New(
		env.AccessExpiry,
		env.AccessSecret,
		env.RefreshExpiry,
		env.RefreshSecret,
	)

	logger.Info("opening a database connection...")

	dbClient, err := config.GetConnection(ctx, env)
	if err != nil {
		return err
	}

	logger.Info("connected to", slog.String("dbName", env.DBName))

	// Services
	userService := user.NewService(cHelper, dbClient, env.DBName)
	authService := auth.NewService(
		userService,
		env.RefreshSecret,
		cHelper,
		logger,
		env.MASTER_EMAIL,
		env.MASTER_PWD,
	)
	articleService := article.NewService(cHelper, dbClient, env.DBName)

	// Controllers
	userController := user.NewController(userService, logger)
	authController := auth.NewController(
		authService,
		logger,
		cHelper,
	)
	articleController := article.NewController(articleService, logger)

	// Middlewares
	jwtVerif := jwtverify.New(
		env.AccessSecret,
		env.RefreshSecret,
		logger,
		cHelper,
	)

	// Server setup
	server := setup.NewServer(
		env,
		logger,
		jwtVerif,
		authController,
		userController,
		articleController,
	)

	// Mux setup
	server.ApplyMiddleWares(mux)

	server.MountHandlers(mux)

	if err := authService.CreateOwnerAccount(ctx); err != nil {
		return err
	}

	logger.Info("listening on", slog.Int64("port", env.Port))

	return nil
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

	w := os.Stderr

	l := slog.New(
		tint.NewHandler(w, &tint.Options{
			TimeFormat: time.Kitchen,
			NoColor:    !isatty.IsTerminal(w.Fd()),
		}),
	)

	env, err := config.LoadEnv()
	if err != nil {
		l.Error("reading from env", slog.String("err", err.Error()))
		return -1
	}

	mux := chi.NewRouter()

	err = startService(ctx, mux, env, l)
	if err != nil {
		l.Error("server setup", slog.String("err", err.Error()))
		return -1
	}

	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", env.Port),
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
		l.Error("shutdown",
			slog.Any("error", err),
		)

		return -1
	}

	l.Info("goodbye 🧩 👋")
	return 0
}
