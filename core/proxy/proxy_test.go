package proxy

import (
	"encoding/json"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/auth-api/core/models"
)

func TestApiError(t *testing.T) {
	err := NewApiError("test error")
	if err.Error() != "test error" {
		t.Fatal("error message should be equal!")
	}
}

func TestCreate(t *testing.T) {
	b := []byte(`{"Username":"wind85","Email":"carlo@email.com","Password":"1234"}`)

	userapi := New()
	buser, _ := userapi.Create(b)
	user := &models.User{}
	err := json.Unmarshal(b, user)
	if err != nil {
		t.Fatal("Unmarshalling user gone wrong!")
	}

	if user.Username != buser.Username {
		t.Fatal("Something is wrong!")
	}
}

func TestUpdate(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.MinCost)

	b := []byte(`{"Username":"carlo85","Email":"carlo@email.com","Password":"` + string(hash) + `"}`)

	userapi := New()
	buser, _ := userapi.Update(b)
	user := &models.User{}
	err := json.Unmarshal(b, user)
	if err != nil {
		t.Fatal("Unmarshalling user gone wrong!")
	}

	if user.Username != buser.Username {
		t.Fatal("Something is wrong!")
	}
}

func TestGet(t *testing.T) {
	b := []byte(`{"Email":"carlo@email.com"}`)

	userapi := New()
	buser, _, _ := userapi.Get(b)
	user := &models.User{}
	err := json.Unmarshal(b, user)
	if err != nil {
		t.Fatal("Unmarshalling user gone wrong!")
	}

	if user.Email != buser.Email {
		t.Fatal("Something is wrong!")
	}
}
