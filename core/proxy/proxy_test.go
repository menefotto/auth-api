package proxy

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
)

func TestApiError(t *testing.T) {
	err := errors.New("test error")
	if err.Error() != "test error" {
		t.Fatal("error message should be equal!")
	}
}

func TestCreate(t *testing.T) {
	u := &models.User{Username: "wind85", Email: "carlo@email.com", Password: "1234"}

	userapi := New()
	buser, _ := userapi.Create(u)
	user := &models.User{}

	if user.Username != buser.Username {
		t.Fatal("Something is wrong!")
	}
}

func TestUpdate(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.MinCost)

	u := &models.User{Username: "wind85", Email: "carlo@email.com", Password: "12345678"}

	userapi := New()
	buser, _ := userapi.Update(u)

	if string(hash) != buser.Username {
		t.Fatal("Something is wrong!")
	}
}

func TestGet(t *testing.T) {
	b := &models.User{Email: "carlo@email.com"}

	userapi := New()
	buser, _ := userapi.Get(b)

	if b.Email != buser.Email {
		t.Fatal("Something is wrong!")
	}
}
