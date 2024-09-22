package asset

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed all:static
var Assets embed.FS

// ServeAssets mount the embedded assets to an HTTP server
func Mount(r chi.Router) {
	r.Route("/static", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		})
		r.Handle("/*", http.FileServer(http.FS(Assets)))
	})
}
