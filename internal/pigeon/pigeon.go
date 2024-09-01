package pigeon

import (
	"fmt"
	"net/smtp"
)

type SMTPRequest struct {
	pigeon Pigeon

	to      []string
	from    string
	subject string
	body    string
}

type Pigeon struct {
	Auth smtp.Auth
	From string
	Addr string
}

func (p Pigeon) NewRequest(to []string, subject, body string) *SMTPRequest {
	return &SMTPRequest{
		pigeon:  p,
		to:      to,
		subject: fmt.Sprintf("Subject: %s", subject),
		body:    body,
		from:    p.From,
	}
}

func (r *SMTPRequest) SendEmail() error {

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(fmt.Sprintf("%s\n", r.subject) + mime + r.body)

	if err := smtp.SendMail(r.pigeon.Addr, r.pigeon.Auth, r.from, r.to, msg); err != nil {
		return err
	}
	return nil
}
