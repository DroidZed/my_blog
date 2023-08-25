package services

import "golang.org/x/crypto/bcrypt"

func HashPassword(txt string) (string, error) {

	result, err := bcrypt.GenerateFromPassword([]byte(txt), 12)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func CompareSecureToPlain(secure string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(secure), []byte(plain)) == nil
}
