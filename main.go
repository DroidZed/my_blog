package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DroidZed/go_lance/config"
	"github.com/DroidZed/go_lance/db"
	"github.com/DroidZed/go_lance/routes"

	"github.com/DroidZed/go_lance/utils"
	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func setupHostWithPort(host string, port int64) string { return fmt.Sprintf("%s:%d", host, port) }

// Entry point, setting up chi and graceful shutdown <3
func main() {

	port, err := config.EnvDbPORT()

	if err != nil {
		log.Fatal("Could not retrieve port!\n")
	}

	addr := setupHostWithPort(config.EnvHost(), port)

	// The HTTP Server
	server := &http.Server{Addr: addr, Handler: service(port)}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Clean up to disconnect
	defer func() {
		if err := db.Client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
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
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	fmt.Println("Goodbye 🧩 👋")
}

func service(port int64) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.CleanPath)
	router.Use(middleware.URLFormat)
	router.Use(httplog.LoggerWithName("CHI API"))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	router.Get("/dev", func(w http.ResponseWriter, r *http.Request) {
		utils.LogAllRoutes(router)
		w.Write([]byte("Nothing will be returned. This is just a dummy message. If you're a developer, check your console."))

	})

	router.Mount("/user", routes.UserRoutes())
	// router.Mount("/dummy", routes.DummyRoute())

	fmt.Printf("Listening on port: %d\n", port)

	return router

}
