package cryptor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DroidZed/my_blog/internal/config"
	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(sub string) (string, error) {

	env := config.LoadEnv()

	daysToAdd, _ := strconv.ParseInt(env.AccessExpiry, 10, 64)

	exp := getExpiration(daysToAdd)

	tokenString, err := createToken(
		sub,
		exp,
		env.AccessSecret,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {

	env := config.LoadEnv()

	daysToAdd, _ := strconv.ParseInt(env.RefreshExpiry, 10, 64)

	exp := getExpiration(daysToAdd)

	v := utils.RNG(17)
	v2 := utils.RNG(19)
	v3 := utils.RNG(149)

	tokenString, err := createToken(
		fmt.Sprintf("0%d1%d2%d", v, v2, v3),
		exp,
		env.RefreshSecret,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getExpiration(daysToAdd int64) int64 {
	duration := time.Duration(daysToAdd) * 24 * time.Hour
	return time.Now().Add(duration).UTC().Unix()
}

func createToken(sub string, expiry int64, sec string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": expiry,
		"iat": time.Now().UTC().Unix(),
	})

	secret := []byte(sec)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractExpiryFromClaims(token *jwt.Token) (float64, error) {

	x, err := extractXFromClaims[float64]("exp", token)

	return *x, err
}

func ExtractSubFromClaims(token *jwt.Token) (string, error) {

	x, err := extractXFromClaims[string]("sub", token)

	return *x, err
}

func extractXFromClaims[T interface{}](claimId string, token *jwt.Token) (*T, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("no claims")
	}

	claim, ok := claims[claimId]
	if !ok {
		return nil, fmt.Errorf("no sub, claims corrupted")
	}

	x, ok := claim.(T)
	if !ok {
		return nil, fmt.Errorf("claim type assertion failed")
	}

	return &x, nil
}

func ParseToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := []byte(secret)

		return secret, nil

	}, jwt.WithValidMethods([]string{"HS256"}))
}
