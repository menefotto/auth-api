package proxy

import (
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/managers"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
)

type UsersVerifier struct {
	mng *managers.Users
}

func New() *UsersVerifier {
	return &UsersVerifier{managers.New("Users")}
}

func (j *UsersVerifier) Create(user *models.User) (*models.User, error) {
	err := j.mng.Verify(user, settings.CREATE_USER_FIELD_REQUIRED)
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
	err := j.mng.Verify(user, settings.UPDATE_USER_FIELD_REQUIRED)
	if err != nil {
		return nil, err
	}

	upUser, err := j.mng.Update(user)
	if err != nil {
		return nil, err
	}

	return upUser, nil
}
