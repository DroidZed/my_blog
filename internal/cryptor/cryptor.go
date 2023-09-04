package cryptor

import (
	"github.com/DroidZed/go_lance/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(txt string) (string, error) {
	bytes := utils.StringToBytes(txt)

	result, err := bcrypt.GenerateFromPassword(bytes, 12)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func CompareSecureToPlain(secure string, plain string) bool {

	secBytes := utils.StringToBytes(secure)
	plainBytes := utils.StringToBytes(plain)

	return bcrypt.CompareHashAndPassword(secBytes, plainBytes) == nil
}
