package httpslog

import (
	"log/slog"
	"net/http"
)

func New(l *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("incoming request",
				slog.String("remoteAddr", r.RemoteAddr),
			)
			next.ServeHTTP(w, r)

		})
	}
}
