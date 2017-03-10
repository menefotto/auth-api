package proxy

import (
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/managers"
	"github.com/auth-api/core/models"
	"github.com/spf13/viper"
)

type UsersVerifier struct {
	mng *managers.Users
}

func New() *UsersVerifier {
	return &UsersVerifier{managers.New("Users")}
}

func (j *UsersVerifier) Create(user *models.User) (*models.User, error) {
	err := j.mng.Verify(user, viper.GetStringSlice("required_user_field.required"))
	if err != nil {
		return nil, err
	}

	newuser, err := j.mng.Create(user)
	if err != nil {
		return nil, errors.InternalDb
	}

	return newuser, nil
}

func (j *UsersVerifier) Get(user *models.User) (*models.User, error) {
	if user.Email == "" {
		return nil, errors.EmailMissing
	}

	gotUser, err := j.mng.Get(user)
	if err != nil {
		return nil, errors.UserNotFound
	}

	return gotUser, nil
}

func (j *UsersVerifier) Update(user *models.User) (*models.User, error) {
	err := j.mng.Verify(user, viper.GetStringSlice("required_user_field.update"))
	if err != nil {
		return nil, err
	}

	upUser, err := j.mng.Update(user)
	if err != nil {
		return nil, err
	}

	return upUser, nil
}
