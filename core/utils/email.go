package utils

import (
	"bytes"
	"errors"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/spf13/viper"
)

type Email struct {
	Title   string
	Message string
}

func SendEmail(sendto []string, msg *Email, templname string) error {
	auth := smtp.PlainAuth(
		"",
		viper.GetString("email.sender"),
		viper.GetString("email.password"),
		viper.GetString("smtp"),
	)

	body, err := RenderEmail(msg, templname, viper.GetString("email.sender"), sendto)
	if err != nil {
		return err
	}

	err = smtp.SendMail(
		viper.GetString("email.smtp")+viper.GetString("email.port"),
		auth,
		viper.GetString("email.sender"),
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
	subj := "Subject: " + viper.GetString("project.id") + ": " + msg.Title + "\r\n\r\n"

	path := filepath.Join(viper.GetString("email.template_dir"), templname+".tmpl")
	tmpl := template.Must(template.ParseFiles(path))

	buff := &bytes.Buffer{}
	err := tmpl.Execute(buff, msg)
	if err != nil {
		return nil, err
	}

	return []byte(send + recv + mime + subj + buff.String() + "\r\n"), nil
}
