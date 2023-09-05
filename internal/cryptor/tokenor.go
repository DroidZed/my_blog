package cryptor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/go-chi/jwtauth"
)

const algorithm = "HS256"

// Set subject ("sub") in the claims
func setSub(claims map[string]interface{}, sub string) {
	claims["sub"] = sub
}

func GenerateAccessToken(sub string) (string, error) {

	daysToAdd, _ := strconv.ParseInt(config.EnvJwtExp(), 10, 64)

	exp := getExpiration(daysToAdd)

	accessClaims := make(map[string]any)

	tokenString, err := createToken(accessClaims, sub, exp)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {

	daysToAdd, _ := strconv.ParseInt(config.EnvRefreshExp(), 10, 64)

	exp := getExpiration(daysToAdd)

	nums := utils.LinearRandomGenerator(89651649874945, 173, 17, 97, 3)

	v := utils.RNG(nums[0])
	v2 := utils.RNG(nums[1])
	v3 := utils.RNG(nums[2])

	refreshClaims := make(map[string]any)

	tokenString, err := createToken(refreshClaims, fmt.Sprintf("0%d1%d2%d", v, v2, v3), exp)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getExpiration(daysToAdd int64) time.Time {

	duration := time.Duration(daysToAdd) * 24 * time.Hour

	return time.Now().Add(duration)
}

func createToken(claims map[string]interface{}, sub string, expiry time.Time) (string, error) {

	tokenAuth := jwtauth.New(algorithm, []byte(config.EnvRefreshSecret()), nil)

	setSub(claims, sub)
	jwtauth.SetExpiry(claims, expiry)
	jwtauth.SetIssuedNow(claims)

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
