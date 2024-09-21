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
	accessKey  string
	refreshKey string
	logger     *slog.Logger
	hasher     cryptor.CryptoHelper
}

func New(
	accessKey string,
	refreshKey string,
	logger *slog.Logger,
	hasher cryptor.CryptoHelper,
) *JwtVerify {
	return &JwtVerify{
		accessKey:  accessKey,
		refreshKey: refreshKey,
		logger:     logger,
		hasher:     hasher,
	}
}

func (j JwtVerify) AccessVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenFromHeader, err := retrieveTokenFromHeader(r)

		if err != nil {
			j.logger.Error("header corrupted", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		token, err := j.validateToken(tokenFromHeader, j.accessKey)

		if err != nil {
			j.logger.Error("invalid token", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		userID, err := j.hasher.ExtractSubFromClaims(token)
		if err != nil {
			j.logger.Error("corrupted sub", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		ctx := context.WithValue(r.Context(), AuthCtxKey{}, userID)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func (j JwtVerify) RefreshVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenFromHeader, err := retrieveTokenFromHeader(r)

		if err != nil {
			j.logger.Error("header corrupted", slog.String("err", err.Error()))
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		_, err = j.validateToken(tokenFromHeader, j.refreshKey)

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

func (j JwtVerify) validateToken(headerValue, secret string) (*jwt.Token, error) {
	token, err := j.hasher.ParseToken(headerValue, secret)
	if err != nil {
		return nil, err
	}

	exp, err := j.hasher.ExtractExpiryFromClaims(token)
	if err != nil {
		return nil, err
	}

	if expired := time.Unix(int64(exp), 0).Before(time.Now()); expired {
		return nil, err
	}

	return token, nil
}
