package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCreateToken(t *testing.T) {
	token := GenerateJwt([]byte("carlolocci"), 3)
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

func TestGenerateCrsf(t *testing.T) {
	token, err := GenerateCrsf("locci.carlo.85@gmail.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(token))
}

func TestDecrypt2(t *testing.T) {
	token := `{"Crsf":"8a5cfc098f003544101bf9c3b47ae1199365159c2a419cd1c3d0ed20de87c08ad5965df4596003169bea37afe0453aecf1b35887"}`

	crsf := &CrsfToken{}

	err := json.Unmarshal([]byte(token), crsf)
	if err != nil {
		t.Fatal(err)
		return
	}

	dec, err := Decrypt(crsf.Token)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("Decrypted:", fmt.Sprintf("%s\n", dec))
}
