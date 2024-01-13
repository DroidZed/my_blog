package middleware

import (
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
)

func FixUrl(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeCtx := chi.RouteContext(r.Context())

		routeP := routeCtx.RoutePath

		if routeP != "" {
			if r.URL.RawPath != "" {
				routeP = r.URL.EscapedPath()
			} else {
				routeP = r.URL.Path
			}

			path.Clean(routeP)

			routeCtx.RoutePath = strings.Replace(routeP, "*", "", strings.Count(routeP, "*"))
		}

		next.ServeHTTP(w, r)
	})
}
