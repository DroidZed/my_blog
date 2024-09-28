package cryptor

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DroidZed/my_blog/internal/utils"
	_ "github.com/joho/godotenv/autoload"

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

var (
	accessKey string = os.Getenv("ACCESS_KEY")
	accessExp string = os.Getenv("ACCESS_EXP")

	refreshKey string = os.Getenv("REFRESH_KEY")
	refreshExp string = os.Getenv("REFRESH_EXP")
)

type Cryptor struct{}

func New() *Cryptor {
	return &Cryptor{}
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

	daysToAdd, _ := strconv.ParseInt(accessExp, 10, 64)

	exp := getExpiration(daysToAdd)

	tokenString, err := createToken(
		exp,
		sub,
		accessKey,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c Cryptor) GenerateRefreshToken() (string, error) {

	daysToAdd, _ := strconv.ParseInt(refreshExp, 10, 64)

	exp := getExpiration(daysToAdd)

	tokenString, err := createToken(
		exp,
		utils.GenUUID(),
		refreshKey,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c Cryptor) ExtractExpiryFromClaims(token *jwt.Token) (int64, error) {

	if !token.Valid {
		return 0, fmt.Errorf("no claims")
	}

	x, err := token.Claims.GetExpirationTime()

	if err != nil {
		return 0, err
	}

	return x.Unix(), nil
}

func (c Cryptor) ExtractSubFromClaims(token *jwt.Token) (string, error) {

	if !token.Valid {
		return "", fmt.Errorf("no claims")
	}

	x, err := token.Claims.GetSubject()

	if err != nil {
		return "", err
	}

	return x, nil
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
