package services

import (
	"encoding/json"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/proxy"
	"github.com/auth-api/core/settings"
	"github.com/auth-api/core/utils"
	"github.com/tcache"
)

type Users struct {
	pool  *proxy.Pool
	cache *tcache.Cache
}

func New(poolsize int) *Users {
	return &Users{
		proxy.NewPool(poolsize),
		tcache.New(time.Minute*8, time.Hour*12),
	}
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

	return utils.GenerateJwt(
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

	url := GenConfirmationUrl(user, "registration", "")

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
	err := u.getUserSendEmail(data, "activation", "activation_confirm", "")
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
		return err
	}

	return nil
}

func (u *Users) PasswordReset(data []byte) error {

	code := utils.GenerateJwt(nil, settings.JWT_PASSWORD_DELTA)
	u.cache.Put(code, data)
	err := u.getUserSendEmail(data, "password reset", "password_reset", code)
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) PasswordResetConfirm(data []byte) error {
	value, ok := u.cache.Get(string(data))
	if !ok {
		return errors.ErrUserNotFound
	}

	code, ok := value.([]byte)
	if !ok {
		return errors.ErrNotValid
	}

	var content interface{}
	err := json.Unmarshal(code, &content)
	if err != nil {
		return errors.ErrJsonPayload
	}

	mapp, ok := content.(map[string]string)
	if !ok {
		return errors.ErrInternalError
	}

	for key, value := range mapp {
		if key == "password" {
			pass, err := bcrypt.GenerateFromPassword(
				[]byte(value), bcrypt.MinCost,
			)
			if err != nil {
				return errors.ErrInternalError
			}
			mng := u.pool.Get()
			defer u.pool.Put(mng)

			msg := `{"` + mapp["email"] + `":"` + string(pass) + `"}`

			user, err := mng.Update(msg)
			if err != nil {
				return err
			}

			err = utils.SendEmail(
				[]string{user.Email},
				&utils.Email{"password reset confirmed", ""},
				"password_reset_confirm",
			)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.ErrFailedPassUpdate
}

func (u *Users) getUserSendEmail(data []byte, title, tmplname, code string) error {
	var url string

	mng := u.pool.Get()
	defer u.pool.Put(mng)

	user, _, err := mng.Get(data)
	if err != nil {
		return err
	}

	if code != "" {
		url = GenConfirmationUrl(user, tmplname, code)
	} else {
		url = GenConfirmationUrl(user, tmplname, user.Code)
	}

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
