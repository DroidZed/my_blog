package pigeon

import (
	"net/smtp"

	"github.com/DroidZed/my_blog/internal/config"
)

var auth smtp.Auth

func GetSmtp() smtp.Auth {

	if auth != nil {
		return auth
	}

	env := config.LoadEnv()

	auth = initAuth(env.SmtpUsername, env.SmtpPassword)

	return auth
}
