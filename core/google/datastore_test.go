package google

import (
	"testing"

	"github.com/auth-api/core/models"
)

// IMPORTANT datastore need the emulator to be running in order to run
// the tests successfully
// ./gcloud beta emulator datastore run
// then export the prompted variable in the env where the test will be run

func TestGstoreOpen(t *testing.T) {
	db := &Datastore{}
	err := db.Open("1345", "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGstoreGetPut(t *testing.T) {
	db := &Datastore{}
	err := db.Open("1345", "test")
	if err != nil {
		t.Fatal(err)
	}
	u := &models.User{}
	u.Username = "wind85"
	err = db.Put("carlo", u)
	if err != nil {
		t.Fatal(err)
	}
	u2, err := db.Get("carlo")
	if err != nil {
		t.Fatal(err)
	}

	if u2.Username != u.Username {
		t.Fatal("Should be identical")
	}
}

func TestGstoreGetPutFail(t *testing.T) {
}

func TestGstoreDel(t *testing.T) {
	db := &Datastore{}
	err := db.Open("1345", "test")
	if err != nil {
		t.Fatal(err)
	}
	u := &models.User{}
	u.Username = "wind85"
	err = db.Put("carlo", u)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Del("carlo")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGstoreNewKey(t *testing.T) {
	db := &Datastore{}
	err := db.Open("1345", "test")
	if err != nil {
		t.Fatal(err)
	}

	k := db.newKey("1234")
	if k == nil {
		t.Fatal(err)
	}

}

func TestGstoreClose(t *testing.T) {
	db := &Datastore{}
	err := db.Open("1345", "test")
	if err != nil {
		t.Fatal(err)
	}

	db.Close()
}

func TestGstoreBackend(t *testing.T) {
	db := &Datastore{}
	err := db.Open("1345", "test")
	if err != nil {
		t.Fatal(err)
	}

	bk := db.Backend()
	if bk == nil {
		t.Fatal(err)
	}

}
