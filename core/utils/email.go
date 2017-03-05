package utils

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"

	"github.com/auth-api/core/settings"
)

func SendEmail(sendto []string, url, templname string) error {
	auth := smtp.PlainAuth(
		"",
		settings.EMAIL_SENDER,
		settings.EMAIL_PASSWORD,
		settings.EMAIL_SMTP,
	)

	body, err := RenderEmail(url, templname)
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

func RenderEmail(url, templname string) ([]byte, error) {
	buff := &bytes.Buffer{}

	path := filepath.Join(settings.EMAIL_TEMPLATE_DIR, templname+".tmpl")
	log.Println("Path: ", path)
	tmpl := template.Must(template.ParseFiles(path))
	err := tmpl.Execute(buff, url)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
