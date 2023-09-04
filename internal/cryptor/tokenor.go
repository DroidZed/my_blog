package cryptor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

const algorithm = "HS256"

func GenerateAccessToken(claims map[string]interface{}) (string, error) {

	daysToAdd, _ := strconv.ParseInt(config.EnvJwtExp(), 10, 64)

	exp, err := createTimestampForToken(daysToAdd, 24*time.Hour)
	if err != nil {
		return "", err
	}

	claims["exp"] = exp

	tokenAuth = jwtauth.New(algorithm, []byte(config.EnvJwtSecret()), nil)

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(claims map[string]interface{}) (string, error) {

	daysToAdd, _ := strconv.ParseInt(config.EnvRefreshExp(), 10, 64)

	exp, err := createTimestampForToken(daysToAdd, 24*time.Hour)
	if err != nil {
		return "", err
	}

	claims["exp"] = exp

	tokenAuth = jwtauth.New(algorithm, []byte(config.EnvRefreshSecret()), nil)

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type numericValue interface {
	float64 | int64 | int32 | float32 | int
}

func createTimestampForToken[ValType numericValue](validity ValType, measurement time.Duration) (int64, error) {

	if validity <= 0 {
		return 0, fmt.Errorf("validity must be greater than zero")
	}

	duration := time.Duration(validity) * measurement

	expDate := time.Now().Add(duration)

	exp := expDate.Unix()

	return exp, nil
}
