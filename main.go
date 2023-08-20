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

	utils "github.com/DroidZed/go_lance/utils"
	"github.com/MadAppGang/httplog"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func setupHostWithPort(host string, port int64) string { return fmt.Sprintf("%s:%d", host, port) }

// Entry point, setting up chi and graceful shutdown <3
func main() {

	port, err := config.EnvDbPORT()

	if err != nil {
		log.Fatal("Could not retrieve port!")
	}

	addr := setupHostWithPort("0.0.0.0", port)

	// The HTTP Server
	server := &http.Server{Addr: addr, Handler: service(port)}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	client := db.GetConnection()

	// Clean up to disconnect
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
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
				log.Fatal("graceful shutdown timed out.. forcing exit.")
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

	fmt.Println("All done, thank you and see you soon 👋")
}

func service(port int64) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(httplog.LoggerWithName("CHI API"))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	router.Get("/dev", func(w http.ResponseWriter, r *http.Request) {
		utils.LogAllRoutes(router)
		w.Write([]byte("Nothing will be returned. This is just a dummy message. If you're a developer, check your console."))

	})

	router.Mount("/user", routes.UserRoutes())

	fmt.Printf("Listening to port: %d", port)

	return router

}
