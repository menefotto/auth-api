package utils

import (
	"bytes"
	"errors"
	"html/template"
	"net/smtp"
	"path/filepath"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"

	"github.com/wind85/auth-api/core/config"
)

type Email struct {
	Title   string
	Message string
}

func SendEmailGun(sendto string, msg *Email, templname string) error {
	sender, err := config.Ini.GetString("mailgun.sender")
	if err != nil {
		return err
	}

	domain, err := config.Ini.GetString("mailgun.domain")
	if err != nil {
		return err
	}

	apikey, err := config.Ini.GetString("mailgun.apikey")
	if err != nil {
		return err
	}

	apipub, err := config.Ini.GetString("mailgun.apipub")
	if err != nil {
		return err
	}

	mg := mailgun.NewMailgun(domain, apikey, apipub)
	body, err := RenderEmailTempl(msg, templname)
	if err != nil {
		return err
	}

	message := mailgun.NewMessage(sender, msg.Title, body, sendto)
	message.AddHeader("Content-Type", "text/html")
	message.AddHeader("MIME-version", "1.0")
	message.SetHtml(body)

	_, _, err = mg.Send(message)
	if err != nil {
		return err
	}

	return err
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

	body, err := RenderEmail(msg, templname, sender, sendto)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", sender, pass, esmtp)
	err = smtp.SendMail(esmtp+port, auth, sender, sendto, body)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil

}

func RenderEmailTempl(msg *Email, templname string) (string, error) {
	tpath, err := config.Ini.GetString("email.template_dir")
	if err != nil {
		return "", err
	}

	path := filepath.Join("../../"+tpath, templname+".tmpl")
	tmpl := template.Must(template.ParseFiles(path))

	buff := &bytes.Buffer{}
	err = tmpl.Execute(buff, msg)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func RenderEmail(msg *Email, templname, from string, sendto []string) ([]byte, error) {
	send := "From: " + from + "\r\n"
	recv := "To: " + sendto[0] + "\r\n"
	mime := "MIME-version: 1.0\r\nContent-Type: text/html\r\n"
	subj := "Subject: " + msg.Title + "\r\n\r\n"

	buff, err := RenderEmailTempl(msg, templname)
	if err != nil {
		return nil, err
	}

	return []byte(send + recv + mime + subj + buff + "\r\n"), nil
}
