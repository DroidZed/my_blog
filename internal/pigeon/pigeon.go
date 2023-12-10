package pigeon

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
)

func ParseHtmlEmail(templateName string) (*template.Template, error) {

	return template.ParseFiles("/public/templates/confirmation_email.html")

}

func NewRequest(to []string, subject, body string) *SMTPRequest {
	return &SMTPRequest{
		to:      to,
		subject: subject,
		body:    body,
		from:    config.LoadEnv().SMTP_USERNAME,
	}
}

func (r *SMTPRequest) SendEmail() (bool, error) {

	smtpAuth := config.GetSmtp()

	env := config.LoadEnv()

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

	msg := utils.StringToBytes(r.subject + mime + "\n" + r.body)

	addr := fmt.Sprintf("%s:%s", env.SMTP_HOST, env.SMTP_PASSWORD)

	if err := smtp.SendMail(addr, smtpAuth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *SMTPRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)

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
