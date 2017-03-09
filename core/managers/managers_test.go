package managers

import (
	"testing"

	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
)

var manager = New("Users")

func TestVerify(t *testing.T) {

	u := &models.User{
		Username: "wind85",
		Password: "1234",
		Email:    "carlo@email.com",
	}

	if err := manager.Verify(u, settings.CREATE_USER_FIELD_REQUIRED); err != nil {
		t.Fatal(err)
	}

	u2 := &models.User{
		Username: "wind85",
		Email:    "carlo@email.com",
	}

	if err := manager.Verify(u2, settings.CREATE_USER_FIELD_REQUIRED); err == nil {
		t.Fatal("should not pass verification")
	}

	u3 := &models.User{
		Username: "wind85",
		Password: "1234",
		Email:    "carlo@email.com",
	}

	err := manager.Verify(u3, settings.UPDATE_USER_FIELD_REQUIRED)
	if err != nil {
		t.Fatal(err)
	}

	u4 := &models.User{
		Username: "wind85",
	}

	if err := manager.Verify(u4, settings.UPDATE_USER_FIELD_REQUIRED); err == nil {
		t.Fatal("should not pass verification")
	}
}

func TestBuildGetCreate(t *testing.T) {
	u := &models.User{
		Username: "wind85",
		Password: "1234",
		Email:    "carlo@email.com",
	}

	uNew, err := manager.Create(u)
	if err != nil {
		t.Fatal(err)
	}

	uNew3, err := manager.Get(&models.User{Email: "carlo@email.com"})
	if err != nil {
		t.Fatal(err)
	}

	if u.Email != uNew.Email || u.Email != uNew3.Email {
		t.Fatal("email should be equal")
	}

	if u.Password != uNew.Password || u.Password != uNew3.Password {
		t.Fatal("email should be equal")
	}

	if u.Username != uNew.Username || u.Username != uNew3.Username {
		t.Fatal("email should be equal")
	}
}

func TestUpdate(t *testing.T) {
	mapped := make(map[string]interface{})
	mapped["Email"] = "carlo@email.com"
	mapped["Username"] = "wind1985"

	u2, err := manager.Update(&models.User{Email: "carlo@email.com", Username: "wind1985"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(mapped)
	if u2.Username == "" {
		t.Fatal("should be equal once updated")
	}
}
