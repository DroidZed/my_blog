package login

import (
	"fmt"

	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/DroidZed/go_lance/internal/user"
)

func ValidateUser(login *LoginBody) (string, error) {

	data, err := user.FindUserByEmail(login.Email)
	if err != nil {
		return "", err
	}

	pwdIsValid := cryptor.CompareSecureToPlain(data.Password, login.Password)

	if !pwdIsValid {
		return "", fmt.Errorf("invalid credentials")
	}

	return data.ID.String(), nil
}

func GenerateLoginTokens(userId string) (string, string, error) {

	access, err := cryptor.GenerateAccessToken(userId)
	if err != nil {

		return "", "", err
	}

	refresh, err := cryptor.GenerateRefreshToken()
	if err != nil {

		return "", "", err
	}

	return access, refresh, nil
}
