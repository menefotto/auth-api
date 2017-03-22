package utils

import (
	"bytes"
	"errors"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/wind85/auth-api/core/config"
)

type Email struct {
	Title   string
	Message string
}

func SendEmail(sendto []string, msg *Email, templname string) error {
	sender, err := config.Ini.GetString("email.sender")
	if err != nil {
		return err
	}

	pass, err := config.Ini.GetString("email.password")
	if err != nil {
		return err
	}

	esmtp, err := config.Ini.GetString("email.smtp")
	if err != nil {
		return err
	}

	port, err := config.Ini.GetString("email.port")
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", sender, pass, esmtp)
	body, err := RenderEmail(msg, templname, sender, sendto)
	if err != nil {
		return err
	}

	err = smtp.SendMail(esmtp+port, auth, sender, sendto, body)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil

}

func RenderEmail(msg *Email, templname, from string, sendto []string) ([]byte, error) {
	send := "From: " + from + "\r\n"
	recv := "To: " + sendto[0] + "\r\n"
	mime := "MIME-version: 1.0\r\nContent-Type: text/html\r\n"
	subj := "Subject: " + msg.Title + "\r\n\r\n"

	tpath, err := config.Ini.GetString("email.template_dir")
	if err != nil {
		return nil, err
	}

	path := filepath.Join(tpath, templname+".tmpl")
	tmpl := template.Must(template.ParseFiles(path))

	buff := &bytes.Buffer{}
	err = tmpl.Execute(buff, msg)
	if err != nil {
		return nil, err
	}

	return []byte(send + recv + mime + subj + buff.String() + "\r\n"), nil
}
