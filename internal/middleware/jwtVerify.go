package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func JwtVerify(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log := config.InitializeLogger().LogHandler

		header := r.Header.Get("Authorization")

		segments := strings.Split(header, " ")
		if len(segments) != 2 || segments[0] != "Bearer" {
			// Invalid header format
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid header format."},
			)
			return
		}
		tokenSegment := segments[1]

		options := jwt.WithValidMethods([]string{"HS256"})

		token, err := jwt.Parse(tokenSegment, func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				utils.JsonResponse(
					w,
					http.StatusUnauthorized,
					utils.DtoResponse{Error: "Invalid signature: Algo tempered."},
				)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			secret := utils.StringToBytes(config.EnvJwtSecret())

			return secret, nil
		}, options)

		if err != nil {
			utils.JsonResponse(
				w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid signature."},
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			log.Error(err)
			utils.JsonResponse(
				w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: err.Error()},
			)
			return
		}

		// test the expiration
		expClaim, ok := claims["exp"]

		if !ok {
			log.Error(err)
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Token expired."},
			)
			return
		}
		exp, ok := expClaim.(float64)
		if !ok {
			// Handle invalid "exp" claim format
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "'exp' claim has an invalid format."},
			)
			return
		}

		expUnixTimestamp := int64(exp)

		expired := time.Unix(expUnixTimestamp, 0).Before(time.Now())

		if expired {
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Token expired."},
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
