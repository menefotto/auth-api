package proxy

import (
	"log"

	"github.com/wind85/auth-api/core/config"
	"github.com/wind85/auth-api/core/errors"
	"github.com/wind85/auth-api/core/managers"
	"github.com/wind85/auth-api/core/models"
)

type Users struct {
	mng *managers.Users
}

func New() *Users {
	backend, err := config.Ini.GetString("database.backend")
	if err != nil {
		log.Println(err)
		return nil
	}

	return &Users{managers.New("Users", backend)}
}

func (j *Users) Create(user *models.User) (*models.User, error) {
	fields, err := config.Ini.GetSlice("required_fields.create")
	if err != nil {
		return nil, err
	}

	err = j.mng.Verify(user, fields)
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
	fields, err := config.Ini.GetSlice("required_fields.create")
	if err != nil {
		return nil, err
	}

	err = j.mng.Verify(user, fields)
	if err != nil {
		return nil, err
	}

	upUser, err := j.mng.Update(user)
	if err != nil {
		return nil, err
	}

	return upUser, nil
}
