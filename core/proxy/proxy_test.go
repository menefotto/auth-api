package proxy

import (
	"testing"

	_ "github.com/wind85/auth-api/core/config"
	"github.com/wind85/auth-api/core/errors"
	"github.com/wind85/auth-api/core/models"
)

func TestApiError(t *testing.T) {
	err := errors.New("test error")
	if err.Error() != "test error" {
		t.Fatal("error message should be equal!")
	}
}

func TestCreate(t *testing.T) {
	u := &models.User{
		Username: "wind85",
		Email:    "carlo@email.com",
		Password: "12345678",
	}

	userapi := New()
	buser, err := userapi.Create(u)
	if err != nil {
		t.Fatal(err)
	}

	if u.Username != buser.Username {
		t.Fatal("Something is wrong!")
	}
}

func TestUpdate(t *testing.T) {
	u := &models.User{Username: "carlo85", Email: "carlo@email.com", Password: "12345678"}

	userapi := New()
	buser, err := userapi.Update(u)
	if err != nil {
		t.Fatal(err)
	}

	if u.Username != buser.Username {
		t.Fatal("Something is wrong!")
	}
}

func TestGet(t *testing.T) {
	b := &models.User{Email: "carlo@email.com"}

	userapi := New()
	buser, err := userapi.Get(b)
	if err != nil {
		t.Fatal(err)
	}

	if b.Email != buser.Email {
		t.Fatal("Something is wrong!")
	}
}
