package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/settings"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type customClaims struct {
	Custom string
	jwt.StandardClaims
}

func GenerateJwt(data []byte, delta int) string {
	mapped := string(data)

	claims := customClaims{
		mapped,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(delta)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "waterandboards",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(GetPrivateKey())
	if err != nil {
		return ""
	}

	return tokenString
}

func ClaimsFromJwt(tok string) (*customClaims, error) {

	token, err := jwt.ParseWithClaims(tok, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.WrongSigningMethod
			}
			// bolocks implementation http://stackoverflow.com/questions/28204385/using-jwt-go-library-key-is-invalid-or-invalid-type
			return GetPubblicKey(), nil
		})
	if err != nil {
		log.Println("claims from jwt")
		return nil, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok && !token.Valid {
		return nil, errors.NotValid
	}

	return claims, nil
}

func computeHmac() []byte {
	h := hmac.New(sha512.New, []byte(settings.HMAC_SECRET))

	return h.Sum(nil)
}

type CrsfToken struct {
	Token string `json:"crsf"`
}

func GenerateCrsf(data string) ([]byte, error) {
	payload, err := Encrypt(data)
	if err != nil {
		return []byte(""), nil
	}

	token := &CrsfToken{payload}

	csrf, err := json.Marshal(token)
	if err != nil {
		return []byte(""), err
	}

	return csrf, nil
}

func Encrypt(data string) (string, error) {
	nonce, _ := hex.DecodeString(settings.NONCE)

	block, err := aes.NewCipher([]byte(settings.CRYPTO_SECRET))
	if err != nil {
		return "", errors.NewCipher
	}

	asecipher, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.NewGCM
	}

	return fmt.Sprintf("%x", asecipher.Seal(nil, nonce, []byte(data), nil)), nil
}

func Decrypt(data string) (string, error) {
	nonce, _ := hex.DecodeString(settings.NONCE)
	text, _ := hex.DecodeString(data)

	block, err := aes.NewCipher([]byte(settings.CRYPTO_SECRET))
	if err != nil {
		return "", errors.NewCipher
	}

	asecipher, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.NewCipher
	}

	dec, err := asecipher.Open(nil, nonce, text, nil)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return fmt.Sprintf("%s\n", dec), nil
}

func RandomGenerator(length int) (string, error) {

	random := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, random); err != nil {
		return "", errors.New(err.Error())
	}

	return fmt.Sprintf("%x\n", random), nil
}

func CheckPassword(p, p2 string) error {

	pass := []byte(p)
	pass2 := []byte(p2)

	err := bcrypt.CompareHashAndPassword(pass, pass2)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New(err.Error())
	}

	if err == bcrypt.ErrHashTooShort {
		return errors.New(err.Error())
	}

	if err != nil {
		return errors.LoginError
	}

	return nil
}

func ValueFromCrsf(crsf string) (string, error) {
	value := &CrsfToken{}

	err := json.Unmarshal([]byte(`{"crsf":"`+crsf+`"}`), value)
	if err != nil {
		return "", err
	}

	email, err := Decrypt(value.Token)
	if err != nil {
		return "", err
	}

	return email[:len(email)-1], nil
}
