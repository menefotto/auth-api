package services

import (
	"testing"

	"github.com/auth-api/core/models"
)

func TestLogin(t *testing.T) {
	user := &models.User{Password: "12345678", Username: "wind1985", Email: "carlo@email.com"}
	service := New(10)
	token, crsf, err := service.Login(user.Email, user.Password)
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
