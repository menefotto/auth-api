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

	body, err := RenderEmail(msg, templname, settings.EMAIL_SENDER, sendto)
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

func RenderEmail(msg *Email, templname, from string, sendto []string) ([]byte, error) {
	send := "From: " + from + "\r\n"
	recv := "To: " + sendto[0] + "\r\n"
	mime := "MIME-version: 1.0\r\nContent-Type: text/html\r\n"
	subj := "Subject: " + settings.PROJECTID + ": " + templname + "\r\n\r\n"

	path := filepath.Join(settings.EMAIL_TEMPLATE_DIR, templname+".tmpl")
	tmpl := template.Must(template.ParseFiles(path))

	buff := &bytes.Buffer{}
	err := tmpl.Execute(buff, msg)
	if err != nil {
		return nil, err
	}

	return []byte(send + recv + mime + subj + buff.String() + "\r\n"), nil
}
