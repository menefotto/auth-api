package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/smtp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/settings"
	jwt "github.com/dgrijalva/jwt-go"
)

var secret = "v97iv7m0mi98BmPoGK81S7sKt1O1UBTV"

type customClaims struct {
	Custom string
	jwt.StandardClaims
}

func GenerateToken(data []byte) string {
	mapped, err := Encrypt(string(data))
	if err != nil {
		return ""
	}

	claims := customClaims{
		mapped,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(settings.JWTExpirationDelta)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "waterandboards",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(computeHmac())
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return tokenString
}

func computeHmac() []byte {
	h := hmac.New(sha512.New, []byte(secret))

	return h.Sum(nil)
}

func Encrypt(data string) (string, error) {
	nonce, _ := hex.DecodeString(settings.NONCE)

	block, err := aes.NewCipher([]byte(settings.CRYPTO_SECRET))
	if err != nil {
		return "", errors.New(err.Error())
	}

	asecipher, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return fmt.Sprintf("%x", asecipher.Seal(nil, nonce, []byte(data), nil)), nil
}

func Decrypt(data string) (string, error) {
	nonce, _ := hex.DecodeString(settings.NONCE)
	text, _ := hex.DecodeString(data)

	block, err := aes.NewCipher([]byte(settings.CRYPTO_SECRET))
	if err != nil {
		return "", errors.New(err.Error())
	}

	asecipher, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New(err.Error())
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
		return errors.ErrLoginError
	}

	return nil
}

func VerifyRequest(cookie string, crsf string) error {

	email, err := Decrypt(crsf)
	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(cookie, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.ErrWrongSigningMethod
			}

			return "", nil
		})

	claims, ok := token.Claims.(*customClaims)
	if !ok && !token.Valid {
		return errors.ErrNotValid
	}

	emailCheck, err := Decrypt(claims.Custom)
	if err != nil {
		return errors.New(err.Error())
	}

	if email != emailCheck {
		return errors.ErrDontMatch
	}

	return nil

}

func SendEmail(sendto []string, body string) error {
	auth := smtp.PlainAuth(
		"",
		settings.EMAIL_SENDER,
		settings.EMAIL_PASSWORD,
		settings.EMAIL_SMTP,
	)

	err := smtp.SendMail(
		settings.EMAIL_SMTP+":"+settings.EMAIL_PORT,
		auth,
		settings.EMAIL_SENDER,
		sendto,
		[]byte(body),
	)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil

}

func GetUserAndEmail(data []byte) error {
	mng := pool.Get()
	defer pool.Put(mng)

	user, err := mng.Create(data)
	if err != nil {
		return err
	}

	err = SendEmail([]string{user.Email}, "default text")
	if err != nil {
		return err
	}

	return nil
}
