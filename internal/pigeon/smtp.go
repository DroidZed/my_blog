package pigeon

import (
	"net/smtp"

	"github.com/DroidZed/go_lance/internal/config"
)

var auth smtp.Auth

func GetSmtp() smtp.Auth {

	if auth != nil {
		return auth
	}

	env := config.LoadEnv()

	auth = initAuth(env.SMTP_USERNAME, env.SMTP_PASSWORD)

	return auth
}
