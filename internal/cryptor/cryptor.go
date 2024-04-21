package cryptor

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(txt string) (string, error) {
	bytes := []byte(txt)

	result, err := bcrypt.GenerateFromPassword(bytes, 12)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func CompareSecureToPlain(secure string, plain string) bool {

	secBytes := []byte(secure)
	plainBytes := []byte(plain)

	return bcrypt.CompareHashAndPassword(secBytes, plainBytes) == nil
}
