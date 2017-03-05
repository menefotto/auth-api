package services

import (
	"testing"
)

func TestLogin(t *testing.T) {
	user := []byte(`{"Password": "12345678", "Username": "wind1985", "Email": "carlo@email.com"}`)
	service := New(10)
	token, crsf, err := service.Login(user)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(token, crsf)
}

/*
func TestSendEmail(t *testing.T) {
	err := SendEmail([]string{"locci.carlo.fb@gmail.com"}, "This is a test email!")
	if err != nil {
		t.Fatal(err)
	}
}
*/
