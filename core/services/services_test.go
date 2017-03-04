package services

import (
	"fmt"
	"testing"
)

func TestCreateToken(t *testing.T) {
	token := GenerateToken([]byte("carlolocci"))
	t.Log("Token", token)
}

func TestHmac(t *testing.T) {
	//t.Log(computeHmac())
}

var enc string

func TestEncryptDecrypt(t *testing.T) {
	enc, _ = Encrypt("ciao carlo")

	_, _ = Encrypt("ciao caro")
	t.Log("Encrypted:", fmt.Sprintf("%s\n", enc))
}

func TestDecrypt(t *testing.T) {
	dec, err := Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Decrypted:", fmt.Sprintf("%s\n", dec))
}

func TestLogin(t *testing.T) {
	user := []byte(`{"Password": "12345678", "Username": "wind1985", "Email": "carlo@email.com"}`)
	token, crsf, err := Login(user)
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
