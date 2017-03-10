package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/auth-api/core/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(data string) (string, error) {
	nonce, _ := hex.DecodeString(viper.GetString("crypto.nonce"))

	block, err := aes.NewCipher([]byte(viper.GetString("crypto.secret")))
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
	nonce, _ := hex.DecodeString(viper.GetString("crypto.nonce"))
	text, _ := hex.DecodeString(data)

	block, err := aes.NewCipher([]byte(viper.GetString("crypto.secret")))
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
