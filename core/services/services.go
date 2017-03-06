package services

import (
	"log"
	"strings"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/proxy"
	"github.com/auth-api/core/settings"
	"github.com/auth-api/core/utils"
)

type Users struct {
	pool *proxy.Pool
}

func New(poolsize int) *Users {
	return &Users{proxy.NewPool(poolsize)}
}

func (u *Users) Login(data []byte) (string, []byte, error) {

	mng := u.pool.Get()
	defer u.pool.Put(mng)

	user, buser, err := mng.Get(data)
	if err != nil {
		return "", nil, err
	}

	err = utils.CheckPassword(
		user.Password,
		buser.Password,
	)
	if err != nil {
		return "", nil, err
	}

	csrf, err := utils.GenerateCrsf(user.Email)
	if err != nil {
		return "", nil, err
	}

	return utils.GenerateToken(
		[]byte(user.Email),
		settings.JWT_LOGIN_DELTA,
	), csrf, nil
}

func (u *Users) Logout(cookie string, crsf string) error {
	err := u.verifyRequest(cookie, crsf)
	if err != nil {
		return err
	}

	// add user blacklisting
	return nil
}

func (u *Users) Me(cookie string, crsf string, data []byte) (*models.User, error) {
	err := u.verifyRequest(cookie, crsf)
	if err != nil {
		return nil, err
	}

	mng := u.pool.Get()
	defer u.pool.Put(mng)

	if data != nil {
		user, err := mng.Update(data)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	email, err := utils.ValueFromCrsf(crsf)
	if err != nil {
		return nil, err
	}

	other, _, err := mng.Get([]byte(`{"email":"` + email + `"}`))
	if err != nil {
		return nil, err
	}

	return other, nil
}

func (u *Users) Register(data []byte) error {
	mng := u.pool.Get()
	defer u.pool.Put(mng)

	user, err := mng.Create(data)
	if err != nil {
		return err
	}

	url := GenConfirmationUrl(user, "registration")

	err = utils.SendEmail(
		[]string{user.Email},
		&utils.Email{"Registration", url},
		"registration",
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) Activation(data []byte) error {
	err := u.getUserSendEmail(data, "activation", "activation_confirm")
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) ActivationConfirm(data []byte) error {
	mng := u.pool.Get()
	defer u.pool.Put(mng)

	claims, err := utils.ClaimsFromJwt(string(data))
	if err != nil {
		return errors.New(err.Error())
	}

	gotUser, _, err := mng.Get([]byte(`{"email":"` + claims.Custom + `"}`))
	if err != nil {
		return err
	}

	if gotUser.Code != string(data) {
		return errors.ErrCodeNotValid
	}

	activatemsg := `{"isactive":"true","email":"` + claims.Custom + `"}`
	user, err := mng.Update([]byte(activatemsg))
	if err != nil {
		return err
	}

	err = utils.SendEmail(
		[]string{user.Email},
		&utils.Email{"Activation Confirmed", ""},
		"activation_confirmation",
	)
	if err != nil {
		log.Println("email err:", err)
		return err
	}

	return nil
}

func (u *Users) PasswordReset(data []byte) error {
	err := u.getUserSendEmail(data, "password reset", "password_reset")
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) PasswordResetConfirm(data []byte) error {

	return nil
}

func (u *Users) getUserSendEmail(data []byte, title, tmplname string) error {
	mng := u.pool.Get()
	defer u.pool.Put(mng)

	user, _, err := mng.Get(data)
	if err != nil {
		return err
	}

	url := GenConfirmationUrl(user, tmplname)

	err = utils.SendEmail(
		[]string{user.Email},
		&utils.Email{title, url},
		tmplname,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) getUserConfirmChange(data []byte, title, tmplname string, update []byte) error {
	mng := u.pool.Get()
	defer u.pool.Put(mng)

	claims, err := utils.ClaimsFromJwt(string(data))
	if err != nil {
		return errors.New(err.Error())
	}

	gotUser, _, err := mng.Get([]byte(`{"email":"` + claims.Custom + `"}`))
	if err != nil {
		return err
	}

	if gotUser.Code != string(data) {
		return errors.ErrCodeNotValid
	}

	active := []byte(`{"isactive":"true","email":"` + claims.Custom + `"}`)
	user, err := mng.Update(active)
	if err != nil {
		return err
	}

	err = utils.SendEmail(
		[]string{user.Email},
		&utils.Email{title, ""},
		tmplname,
	)
	if err != nil {
		return err
	}

	return nil

}

func (u *Users) verifyRequest(cookie string, crsf string) error {
	email, err := utils.ValueFromCrsf(crsf)
	if err != nil {
		//log.Println("here")
		return err
	}
	claims, err := utils.ClaimsFromJwt(cookie)
	if err != nil {
		return err
	}
	if strings.Compare(email, claims.Custom) != 0 {
		return errors.ErrDontMatch
	}

	return nil
}
