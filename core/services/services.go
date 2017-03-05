package services

import (
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/proxy"
)

var pool = proxy.NewPool(10)

func Login(data []byte) (string, []byte, error) {

	mng := pool.Get()
	defer pool.Put(mng)

	user, buser, err := mng.Get(data)
	if err != nil {
		return "", []byte(""), err
	}

	err = CheckPassword(user.Password, buser.Password)
	if err != nil {
		return "", []byte(""), err
	}

	csrf, err := GenerateCrsf(user.Email)
	if err != nil {
		return "", []byte(""), err
	}

	return GenerateToken([]byte(user.Email)), csrf, nil
}

func Logout(cookie string, crsf string) error {
	err := VerifyRequest(cookie, crsf)
	if err != nil {
		return err
	}

	// add user blacklisting
	return nil
}

func Me(cookie string, crsf string, data []byte) (*models.User, error) {
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

	email, err := ValueFromCrsf(crsf)
	if err != nil {
		return nil, err
	}

	other, _, err := mng.Get([]byte(`{"email":"` + email + `"}`))
	if err != nil {
		return nil, err
	}

	return other, nil
}

func Registration(data []byte) error {
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

func Activate(data []byte) error {
	mng := pool.Get()
	defer pool.Put(mng)

	user, err := mng.Update(data)
	if err != nil {
		return err
	}

	err = SendEmail([]string{user.Email}, "default text")
	if err != nil {
		return err
	}

	return nil
}

func PasswordReset(data []byte) error {
	err := GetUserAndEmail(data)
	if err != nil {
		return err
	}

	return nil
}

func PasswordResetConfirm(data []byte) error {
	err := GetUserAndEmail(data)
	if err != nil {
		return err
	}

	return nil
}
