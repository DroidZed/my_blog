package middleware

import (
	"net/http"
	"path"
)

func FixUrl(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.RequestURI = path.Clean(r.URL.EscapedPath())

		next.ServeHTTP(w, r)
	})
}
