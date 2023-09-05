package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func JwtVerify(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log := config.InitializeLogger().LogHandler

		tokenString := jwtauth.TokenFromHeader(r)

		// test the expiration

		tokenBytes := utils.StringToBytes(tokenString)

		validationOpts := jwt.WithVerify(jwa.HS256, config.EnvJwtSecret())

		jwtObj, err := jwt.Parse(tokenBytes, validationOpts) //validationOpts
		if err != nil {
			log.Error(err)
			http.Error(w, fmt.Errorf("invalid token. Could not parse token").Error(), http.StatusUnauthorized)
			return
		}

		exp := jwtObj.Expiration()

		if time.Now().After(exp) {
			http.Error(w, fmt.Errorf("invalid token. Expired").Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
