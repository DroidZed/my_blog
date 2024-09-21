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
	HashPlain(txt string) (string, error)
	CompareSecureToPlain(secure string, plain string) bool
	GenerateAccessToken(sub string) (string, error)
	GenerateRefreshToken() (string, error)
	ExtractExpiryFromClaims(token *jwt.Token) (int64, error)
	ExtractSubFromClaims(token *jwt.Token) (string, error)
	ParseToken(token string, secret string) (*jwt.Token, error)
}

type Cryptor struct {
	accessExp  string
	accessKey  string
	refreshExp string
	refreshKey string
}

func New(
	accessExp string,
	accessKey string,
	refreshExp string,
	refreshKey string,
) *Cryptor {
	return &Cryptor{
		accessExp:  accessExp,
		accessKey:  accessKey,
		refreshExp: refreshExp,
		refreshKey: refreshKey,
	}
}

func (c Cryptor) HashPlain(txt string) (string, error) {
	bytes := []byte(txt)

	result, err := bcrypt.GenerateFromPassword(bytes, 12)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (c Cryptor) CompareSecureToPlain(secure string, plain string) bool {

	secBytes := []byte(secure)
	plainBytes := []byte(plain)

	return bcrypt.CompareHashAndPassword(secBytes, plainBytes) == nil
}

func (c Cryptor) GenerateAccessToken(sub string) (string, error) {

	daysToAdd, _ := strconv.ParseInt(c.accessExp, 10, 64)

	exp := getExpiration(daysToAdd)

	tokenString, err := createToken(
		exp,
		sub,
		c.accessKey,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c Cryptor) GenerateRefreshToken() (string, error) {

	daysToAdd, _ := strconv.ParseInt(c.refreshExp, 10, 64)

	exp := getExpiration(daysToAdd)

	tokenString, err := createToken(
		exp,
		utils.GenUUID(),
		c.refreshKey,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c Cryptor) ExtractExpiryFromClaims(token *jwt.Token) (int64, error) {

	x, err := extractXFromClaims[int64]("exp", token)

	return *x, err
}

func (c Cryptor) ExtractSubFromClaims(token *jwt.Token) (string, error) {

	x, err := extractXFromClaims[string]("sub", token)

	return *x, err
}

func (c Cryptor) ParseToken(token string, secret string) (*jwt.Token, error) {
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
