package amazon

import (
	"testing"

	"github.com/auth-api/core/models"
)

func TestPut(t *testing.T) {
	d := &Dynamo{}
	err := d.Open("Users", "us-west-2")
	if err != nil {
		t.Fatal(err)
	}

	u := &models.User{Email: "carlo@example.com", Password: "123456"}
	err = d.Put(u.Email, u)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	d := &Dynamo{}
	err := d.Open("Users", "us-west-2")
	if err != nil {
		t.Fatal(err)
	}

	u := &models.User{Email: "carlo@example.com", Password: "123456"}
	u2, err := d.Get(u.Email)
	if err != nil {
		t.Fatal(err)
	}

	if u.Email != u2.Email {
		t.Fatal("Emails should be equal")
	}
}

func TestDel(t *testing.T) {
	d := &Dynamo{}
	err := d.Open("Users", "us-west-2")
	if err != nil {
		t.Fatal(err)
	}

	err = d.Del("carlo@example.com")
	if err != nil {
		t.Fatal("Delete error: ", err)
	}
}

func TestSatisfyInterface(t *testing.T) {
}
