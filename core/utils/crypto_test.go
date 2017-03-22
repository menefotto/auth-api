package utils

import (
	"fmt"
	"testing"

	_ "github.com/wind85/auth-api/core/config"
)

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

func TestDecrypt2(t *testing.T) {
	test := "8a5cfc098f003544101bf9c3b47ae1199365159c2a419cd1c3d0ed20de87c08ad5965df4596003169bea37afe0453aecf1b35887"

	dec, err := Decrypt(test)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("Decrypted:", fmt.Sprintf("%s\n", dec))
}
