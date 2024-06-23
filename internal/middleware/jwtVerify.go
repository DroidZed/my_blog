package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthCtxKey struct{}

func AccessVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := config.LoadEnv()
		log := config.GetLogger()

		tokenFromHeader, err := retrieveTokenFromHeader(r)

		if err != nil {
			log.Errorf("Header corrupted.\n%s", err)
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		token, err := validateToken(tokenFromHeader, env.RefreshSecret)

		if err != nil {
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		userId, err := cryptor.ExtractSubFromClaims(token)
		if err != nil {
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "No sub in token."},
			)
			return
		}

		ctx := context.WithValue(r.Context(), AuthCtxKey{}, userId)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func RefreshVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := config.LoadEnv()
		log := config.GetLogger()

		tokenFromHeader, err := retrieveTokenFromHeader(r)

		if err != nil {
			log.Errorf("Header corrupted.\n%s", err)
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		_, err = validateToken(tokenFromHeader, env.RefreshSecret)

		if err != nil {
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func retrieveTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	segments := strings.Split(header, " ")
	if len(segments) != 2 || segments[0] != "Bearer" {
		return "", fmt.Errorf("invalid format")
	}

	return segments[1], nil
}

func validateToken(headerValue, secret string) (*jwt.Token, error) {
	log := config.GetLogger()

	token, err := cryptor.ParseToken(headerValue, secret)
	if err != nil {
		log.Errorf("Unable to parse token.\n%s", err)
		return nil, err
	}

	exp, err := cryptor.ExtractExpiryFromClaims(token)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if expired := time.Unix(int64(exp), 0).Before(time.Now()); expired {
		return nil, err
	}

	return token, nil
}
