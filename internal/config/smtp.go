package config

import "net/smtp"

var auth smtp.Auth

func GetSmtp() smtp.Auth {

	if auth != nil {
		return auth
	}

	env := LoadEnv()

	auth = smtp.PlainAuth("", env.SMTP_USERNAME, env.SMTP_PASSWORD, env.SMTP_HOST)

	return auth
}
