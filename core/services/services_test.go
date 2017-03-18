package services

import (
	"testing"

	_ "github.com/auth-api/core/config"
	"github.com/auth-api/core/models"
)

func TestLogin(t *testing.T) {
	user := &models.User{
		Password: "12345678",
		Username: "wind1985",
		Email:    "carlo@email.com",
	}

	service := New(10)

	token, crsf, err := service.Login(user.Email, user.Password)
	if err != nil {
		t.Fatal(err)
	}

	if len(token) < 12 && len(string(crsf)) < 12 {
		t.Fatal("Token and crsf too short something is wrong!")
	}

}

/*
func TestSendEmail(t *testing.T) {
	err := SendEmail([]string{"locci.carlo.fb@gmail.com"}, "This is a test email!")
	if err != nil {
		t.Fatal(err)
	}
}
*/
