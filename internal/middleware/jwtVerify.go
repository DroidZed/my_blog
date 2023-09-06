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

		tokenString, err := retrieveTokenFromHeader(r)
		if err != nil {
			log.Errorf("Header corrupted.\n%s", err)
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
		}

		token, err := parseToken(w, tokenString)
		if err != nil {
			log.Errorf("Unable to parse token.\n%s", err)
			utils.JsonResponse(
				w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		exp, err := extractExpiryFromClaims(token)
		if err != nil {
			log.Error(err)
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
		}

		expired := time.Unix(int64(exp), 0).Before(time.Now())
		if expired {
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractExpiryFromClaims(token *jwt.Token) (float64, error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {

		return 0, fmt.Errorf("no claims")
	}

	// test the expiration
	expClaim, ok := claims["exp"]
	if !ok {
		return 0, fmt.Errorf("no expiration, claims corrupted")
	}

	return expClaim.(float64), nil

}

func retrieveTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	segments := strings.Split(header, " ")

	if len(segments) != 2 || segments[0] != "Bearer" {
		return "", fmt.Errorf("invalid format")
	}

	return segments[1], nil
}

func parseToken(w http.ResponseWriter, token string) (*jwt.Token, error) {

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := utils.StringToBytes(config.EnvJwtSecret())

		return secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
}
