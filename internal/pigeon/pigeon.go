package pigeon

import (
	"fmt"
	"net"
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
	smtpAuth smtp.Auth
	from     string
	addr     string
}

func New(
	smtpUsername,
	smtpPwd,
	smtpHost,
	smtpPort string,
) *Pigeon {
	return &Pigeon{
		smtpAuth: smtp.PlainAuth(
			"pigeon",
			smtpUsername,
			smtpPwd,
			smtpHost,
		),
		from: smtpUsername,
		addr: net.JoinHostPort(smtpHost, smtpPort),
	}
}

func (p Pigeon) NewRequest(to []string, subject, body string) *SMTPRequest {
	return &SMTPRequest{
		pigeon:  p,
		to:      to,
		subject: fmt.Sprintf("Subject: %s", subject),
		body:    body,
		from:    p.from,
	}
}

func (r *SMTPRequest) SendEmail() error {

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(fmt.Sprintf("%s\n", r.subject) + mime + r.body)

	if err := smtp.SendMail(r.pigeon.addr, r.pigeon.smtpAuth, r.from, r.to, msg); err != nil {
		return err
	}
	return nil
}
