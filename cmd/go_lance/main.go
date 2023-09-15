package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DroidZed/go_lance/internal"
	"github.com/DroidZed/go_lance/internal/config"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func init() {
	// Register a directive named "path" to retrieve values from `chi.URLParam`,
	// i.e. decode path variables.
	httpin.UseGochiURLParam("path", chi.URLParam)
}

func startService() *internal.Server {

	log := config.InitializeLogger().LogHandler

	server := internal.CreateNewServer()

	envPort := server.EnvConfig.Port

	server.ApplyMiddleWares()

	server.MountHandlers()

	log.Infof("Listening on port: %d\n", envPort)

	return server
}

// Entry point, setting up chi and graceful shutdown <3
// @title GoLance API Docs
// @version 1.0
// @description This is the GoLance API documentation.
// @termsOfService http://example.com/terms/

// @contact.name GoLance Support
// @contact.url http://example.com/support
// @contact.email joe@example.com

// @license.name MIT
// @license.url https://github.com/DroidZed/go_lance/LICENSE

// @host golance.io
// @BasePath /
func main() {

	log := config.InitializeLogger().LogHandler

	app := startService()

	dbClient := app.DbClient
	envPort := app.EnvConfig.Port

	// The HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", envPort),
		Handler: app.Router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Clean up to disconnect
	defer func() {
		if err := dbClient.Disconnect(serverCtx); err != nil {
			log.Fatal(err)
		}
	}()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 3 seconds
		shutdownCtx, cancelFunc := context.WithTimeout(serverCtx, 3*time.Second)
		defer cancelFunc()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {

				log.Fatal("Graceful shutdown timed out.. forcing exit.\n")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	log.Info("Goodbye 🧩 👋")
}
