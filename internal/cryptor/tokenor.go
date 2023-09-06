package cryptor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(sub string) (string, error) {

	daysToAdd, _ := strconv.ParseInt(config.EnvJwtExp(), 10, 64)

	exp := getExpiration(daysToAdd)

	accessClaims := make(map[string]interface{})

	tokenString, err := createToken(
		accessClaims,
		sub,
		exp,
		config.EnvJwtSecret(),
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {

	daysToAdd, _ := strconv.ParseInt(config.EnvRefreshExp(), 10, 64)

	exp := getExpiration(daysToAdd)

	v := utils.RNG(17)
	v2 := utils.RNG(19)
	v3 := utils.RNG(149)

	refreshClaims := make(map[string]interface{})

	tokenString, err := createToken(
		refreshClaims,
		fmt.Sprintf("0%d1%d2%d", v, v2, v3),
		exp,
		config.EnvRefreshSecret(),
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

func createToken(claims map[string]interface{}, sub string, expiry int64, sec string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": expiry,
		"iat": time.Now().UTC().Unix(),
	})

	secret := utils.StringToBytes(sec)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
