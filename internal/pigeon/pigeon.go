package pigeon

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
)

func NewRequest(to []string, subject, body string) *SMTPRequest {
	return &SMTPRequest{
		to:      to,
		subject: subject,
		body:    body,
		from:    config.LoadEnv().SMTP_USERNAME,
	}
}

func (r *SMTPRequest) GetBody() string {
	return r.body
}

func (r *SMTPRequest) SendEmail() error {

	smtpAuth := GetSmtp()

	env := config.LoadEnv()

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := utils.StringToBytes(r.subject + mime + "\n" + r.body)

	addr := fmt.Sprintf("%s:%s", env.SMTP_HOST, env.SMTP_PORT)

	if err := smtp.SendMail(addr, smtpAuth, r.from, r.to, msg); err != nil {
		return err
	}
	return nil
}

func (r *SMTPRequest) ParseTemplate(templateFileBaseName string, data interface{}) error {
	fullName := fmt.Sprintf("%s.tmpl", templateFileBaseName)
	t, err := template.New(fullName).ParseFiles(fmt.Sprintf("public/templates/%s", fullName))

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, data); err != nil {
		return err
	}

	r.body = buf.String()

	return nil
}
