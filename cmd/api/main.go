package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DroidZed/my_blog/internal/article"
	"github.com/DroidZed/my_blog/internal/auth"
	"github.com/DroidZed/my_blog/internal/cryptor"
	_ "github.com/joho/godotenv/autoload"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/DroidZed/my_blog/internal/database"
	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/setup"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

var (
	port string = os.Getenv("PORT")
	host string = os.Getenv("HOST")
)

func startService(
	ctx context.Context,
	mux *chi.Mux,
	logger *slog.Logger,
) error {

	logger.Info("opening a database connection...")

	db, err := database.New(ctx)
	if err != nil {
		return err
	}

	logger.Info("connected to", slog.String("dbName", db.Name))

	// pg := pigeon.New()

	var markdown = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithHardWraps(),
		),
	)

	hasher := cryptor.New()

	// Middlewares
	jwtVerif := jwtverify.New(logger, hasher)

	// Services
	articleService := article.NewService(db, markdown)
	userService := user.NewService(hasher, db)
	authService := auth.NewService(userService, hasher, logger)

	// Controllers
	articleController := article.NewController(articleService, logger, userService)
	authController := auth.NewController(authService, logger, hasher)
	userController := user.NewController(userService, logger)

	// Server setup
	server := setup.NewServer(
		logger,
		articleController,
		userController,
		authController,
		jwtVerif,
		db,
	)

	// Mux setup
	server.ApplyMiddleWares(mux)

	server.MountHandlers(mux)

	authService.CreateOwnerAccount(ctx)

	logger.Info("listening on", slog.String("port", port))

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

	mux := chi.NewRouter()

	err := startService(ctx, mux, l)
	if err != nil {
		l.Error("server setup", slog.String("err", err.Error()))
		return -1
	}

	server := &http.Server{
		Addr:     net.JoinHostPort(host, port),
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
