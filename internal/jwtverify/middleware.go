package jwtverify

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthCtxKey struct{}

type JwtVerify struct {
	AccessSecret  string
	RefreshSecret string
	Logger        *slog.Logger
	CHelper       cryptor.CryptoHelper
}

func (j JwtVerify) AccessVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenFromHeader, err := retrieveTokenFromHeader(r)

		if err != nil {
			j.Logger.Error("header corrupted", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		token, err := validateToken(
			tokenFromHeader,
			j.AccessSecret,
			j.CHelper.ParseToken,
			j.CHelper.ExtractExpiryFromClaims,
		)

		if err != nil {
			j.Logger.Error("invalid token", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		userId, err := j.CHelper.ExtractSubFromClaims(token)
		if err != nil {
			j.Logger.Error("corrupted sub", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		ctx := context.WithValue(r.Context(), AuthCtxKey{}, userId)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func (j JwtVerify) RefreshVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenFromHeader, err := retrieveTokenFromHeader(r)

		if err != nil {
			j.Logger.Error("header corrupted", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		_, err = validateToken(
			tokenFromHeader,
			j.RefreshSecret,
			j.CHelper.ParseToken,
			j.CHelper.ExtractExpiryFromClaims,
		)

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

func validateToken(
	headerValue, secret string,
	ParseToken func(token string, secret string) (*jwt.Token, error),
	ExtractExpiryFromClaims func(token *jwt.Token) (int64, error),
) (*jwt.Token, error) {
	token, err := ParseToken(headerValue, secret)
	if err != nil {
		return nil, err
	}

	exp, err := ExtractExpiryFromClaims(token)
	if err != nil {
		return nil, err
	}

	if expired := time.Unix(int64(exp), 0).Before(time.Now()); expired {
		return nil, err
	}

	return token, nil
}
