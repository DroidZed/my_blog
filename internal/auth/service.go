package auth

import (
	"fmt"

	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/DroidZed/go_lance/internal/user"
)

func ValidateUser(login *LoginBody) (*user.User, error) {

	userService := &user.UserService{}

	data, err := userService.FindUserByEmail(login.Email)
	if err != nil {
		return nil, err
	}

	pwdIsValid := cryptor.CompareSecureToPlain(data.Password, login.Password)

	if !pwdIsValid {
		return nil, fmt.Errorf("invalid credentials")
	}

	return data, nil
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
