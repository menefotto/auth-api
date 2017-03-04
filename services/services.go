package services

import (
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/proxy"
)

var pool = proxy.NewPool(10)

func Login(data []byte) (string, string, error) {

	mng := pool.Get()
	defer pool.Put(mng)

	user, buser, err := mng.Get(data)
	if err != nil {
		return "", "", err
	}

	err = CheckPassword(user.Password, buser.Password)
	if err != nil {
		return "", "", err
	}

	uuid, err := Encrypt(user.Email)
	if err != nil {
		return "", "", nil
	}

	return GenerateToken([]byte(user.Uuid)), uuid, nil
}

func Logout(cookie map[string]string, crsf string) error {
	err := VerifyRequest(cookie, crsf)
	if err != nil {
		return err
	}

	// add user blacklisting
	return nil
}

func Me(cookie map[string]string, crsf string, data []byte) (*models.User, error) {
	err := VerifyRequest(cookie, crsf)
	if err != nil {
		return nil, err
	}

	mng := pool.Get()
	defer pool.Put(mng)

	if data != nil {
		_, err := mng.Update(data)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, user, err := mng.Get(data)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Registration(data []byte) (*models.User, error) {
	mng := pool.Get()
	defer pool.Put(mng)

	user, err := mng.Create(data)
	if err != nil {
		return nil, err
	}

	err := SendEmail(user.Email, "default text")
	if err != nil {
		return user, err
	}

	return user, nil
}

func Activate() {

}

func PasswordReset(data []byte) error {
	mng := pool.Get()
	defer pool.Put(mng)

	user, err := mng.Get(data)
	if err != nil {
		return err
	}

	if err := SendEmail(user.Email, "text body"); err != nil {
		return err
	}

	return nil
}

func PasswordResetConfirm() {
}
