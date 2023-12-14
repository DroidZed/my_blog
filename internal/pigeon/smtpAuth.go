package pigeon

import (
	"errors"
	"net/smtp"
)

type smtpAuth struct {
	username, password string
}

func initAuth(username, password string) smtp.Auth {
	return &smtpAuth{username, password}
}

func (a *smtpAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *smtpAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown fromServer")
		}
	}
	return nil, nil
}
