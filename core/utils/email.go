package utils

import (
	"bytes"
	"errors"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/auth-api/core/settings"
)

type Email struct {
	Title   string
	Message string
}

func SendEmail(sendto []string, msg *Email, templname string) error {
	auth := smtp.PlainAuth(
		"",
		settings.EMAIL_SENDER,
		settings.EMAIL_PASSWORD,
		settings.EMAIL_SMTP,
	)

	body, err := RenderEmail(msg, templname)
	if err != nil {
		return err
	}

	err = smtp.SendMail(
		settings.EMAIL_SMTP+":"+settings.EMAIL_PORT,
		auth,
		settings.EMAIL_SENDER,
		sendto,
		body,
	)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil

}

func RenderEmail(msg *Email, templname string) ([]byte, error) {
	buff := &bytes.Buffer{}

	path := filepath.Join(settings.EMAIL_TEMPLATE_DIR, templname+".tmpl")
	tmpl := template.Must(template.ParseFiles(path))
	err := tmpl.Execute(buff, msg)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
