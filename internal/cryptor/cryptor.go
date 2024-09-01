package cryptor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CryptoHelper interface {
	HashPassword(txt string) (string, error)
	CompareSecureToPlain(secure string, plain string) bool
	GenerateAccessToken(sub string) (string, error)
	GenerateRefreshToken() (string, error)
	ExtractExpiryFromClaims(token *jwt.Token) (int64, error)
	ExtractSubFromClaims(token *jwt.Token) (string, error)
	ParseToken(token string, secret string) (*jwt.Token, error)
}

type Cryptor struct {
	AccessExpiry  string
	AccessSecret  string
	RefreshExpiry string
	RefreshSecret string
}

func (Cryptor) HashPassword(txt string) (string, error) {
	bytes := []byte(txt)

	result, err := bcrypt.GenerateFromPassword(bytes, 12)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (Cryptor) CompareSecureToPlain(secure string, plain string) bool {

	secBytes := []byte(secure)
	plainBytes := []byte(plain)

	return bcrypt.CompareHashAndPassword(secBytes, plainBytes) == nil
}

func (t *Cryptor) GenerateAccessToken(sub string) (string, error) {

	daysToAdd, _ := strconv.ParseInt(t.AccessExpiry, 10, 64)

	exp := getExpiration(daysToAdd)

	tokenString, err := createToken(
		exp,
		sub,
		t.AccessSecret,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Cryptor) GenerateRefreshToken() (string, error) {

	daysToAdd, _ := strconv.ParseInt(t.RefreshExpiry, 10, 64)

	exp := getExpiration(daysToAdd)

	tokenString, err := createToken(
		exp,
		utils.GenUUID(),
		t.RefreshSecret,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Cryptor) ExtractExpiryFromClaims(token *jwt.Token) (int64, error) {

	x, err := extractXFromClaims[int64]("exp", token)

	return *x, err
}

func (t *Cryptor) ExtractSubFromClaims(token *jwt.Token) (string, error) {

	x, err := extractXFromClaims[string]("sub", token)

	return *x, err
}

func (t *Cryptor) ParseToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil

	}, jwt.WithValidMethods([]string{"HS256"}))
}

func getExpiration(daysToAdd int64) int64 {
	duration := time.Duration(daysToAdd) * 24 * time.Hour
	return time.Now().Add(duration).UTC().Unix()
}

func createToken(expiry int64, sub, sec string) (string, error) {

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
