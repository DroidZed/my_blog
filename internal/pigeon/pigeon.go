package pigeon

import (
	"fmt"
	"net/smtp"

	"github.com/DroidZed/my_blog/internal/config"
)

type SMTPRequest struct {
	to      []string
	from    string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *SMTPRequest {
	return &SMTPRequest{
		to:      to,
		subject: fmt.Sprintf("Subject: %s", subject),
		body:    body,
		from:    config.LoadEnv().SmtpUsername,
	}
}

func (r *SMTPRequest) GetBody() string {
	return r.body
}

func (r *SMTPRequest) GetSubject() string {
	return r.subject
}

func (r *SMTPRequest) SendEmail() error {

	smtpAuth := GetSmtp()

	env := config.LoadEnv()

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(fmt.Sprintf("%s\n", r.GetSubject()) + mime + r.GetBody())

	addr := fmt.Sprintf("%s:%s", env.SmtpHost, env.SmtpPort)

	if err := smtp.SendMail(addr, smtpAuth, r.from, r.to, msg); err != nil {
		return err
	}
	return nil
}

func (r *SMTPRequest) SetBody(str string) {
	r.body = str
}
