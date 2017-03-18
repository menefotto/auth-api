package managers

import (
	"testing"

	_ "github.com/auth-api/core/config"
	"github.com/auth-api/core/models"
	"github.com/spf13/viper"
)

var manager = New("Users", "DATASTORE")

func TestVerify(t *testing.T) {

	u := &models.User{
		Username: "wind85",
		Password: "1234",
		Email:    "carlo@email.com",
	}
	if err := manager.Verify(u, viper.GetStringSlice("required_user_fields.create")); err != nil {
		t.Fatal(err)
	}

	u2 := &models.User{
		Username: "wind85",
		Email:    "carlo@email.com",
	}

	if err := manager.Verify(u2, viper.GetStringSlice("required_fields.create")); err == nil {
		t.Fatal("should not pass verification")
	}

	u3 := &models.User{
		Username: "wind85",
		Password: "1234",
		Email:    "carlo@email.com",
	}

	err := manager.Verify(u3, viper.GetStringSlice("required_fields.update"))
	if err != nil {
		t.Fatal(err)
	}

	u4 := &models.User{
		Username: "wind85",
	}

	if err := manager.Verify(u4, viper.GetStringSlice("required_fields.update")); err == nil {
		t.Fatal("should not pass verification")
	}
}

func TestBuildGetCreate(t *testing.T) {
	u := &models.User{
		Username: "wind85",
		Password: "12345678",
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

func TestSetField(t *testing.T) {
	u := &models.User{}

	field1 := "Datejoined"
	field2 := "Password"
	value := "12345678"

	err := SetField(u, field1, value)
	if err == nil {
		t.Fatal("Ops should give an error!")
	}
	err = SetField(u, field2, value)
	if err != nil {
		t.Fatal("Ops should not return this:", err)
	}

}
