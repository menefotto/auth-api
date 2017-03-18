package proxy

import (
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/managers"
	"github.com/auth-api/core/models"
	"github.com/spf13/viper"
)

type Users struct {
	mng *managers.Users
}

func New() *Users {
	return &Users{
		managers.New(
			"Users",
			viper.GetString("database.backend"),
		),
	}
}

func (j *Users) Create(user *models.User) (*models.User, error) {
	err := j.mng.Verify(user, viper.GetStringSlice("required_fields.create"))
	if err != nil {
		return nil, err
	}

	newuser, err := j.mng.Create(user)
	if err != nil {
		return nil, errors.InternalDb
	}

	return newuser, nil
}

func (j *Users) Get(user *models.User) (*models.User, error) {
	if user.Email == "" {
		return nil, errors.EmailMissing
	}

	gotUser, err := j.mng.Get(user)
	if err != nil {
		return nil, errors.UserNotFound
	}

	return gotUser, nil
}

func (j *Users) Update(user *models.User) (*models.User, error) {
	err := j.mng.Verify(user, viper.GetStringSlice("required_fields.update"))
	if err != nil {
		return nil, err
	}

	upUser, err := j.mng.Update(user)
	if err != nil {
		return nil, err
	}

	return upUser, nil
}
