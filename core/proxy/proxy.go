package proxy

import (
	"encoding/json"

	"github.com/auth-api/core/managers"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
)

type UsersJson struct {
	mng *managers.Users
}

func New() *UsersJson {
	return &UsersJson{managers.New("Users")}
}

func (j *UsersJson) Create(data []byte) (*models.User, error) {
	user, err := j.Build(data, "CREATE")
	if err != nil {
		return nil, err
	}

	newuser, err := j.mng.Create(user)
	if err != nil {
		return nil, ErrInternalDb
	}

	return newuser, nil
}

func (j *UsersJson) Get(data []byte) (*models.User, *models.User, error) {
	user, err := j.Build(data, "GET")
	if err != nil {
		return nil, nil, err
	}

	gotUser, err := j.mng.Get(user.Email)
	if err != nil {
		return nil, nil, ErrUserNotFound
	}

	return gotUser, user, nil
}

func (j *UsersJson) Update(data []byte) (*models.User, error) {
	var objects interface{}

	err := json.Unmarshal(data, &objects)
	if err != nil {
		return nil, ErrMalformedInput
	}

	mapped, ok := objects.(map[string]interface{})
	if !ok {
		return nil, ErrMalformedInput
	}

	user, err := j.mng.Update(mapped)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (j *UsersJson) Build(data []byte, t string) (*models.User, error) {
	user := &models.User{}

	err := json.Unmarshal(data, user)
	if err != nil {
		return nil, ErrMalformedInput
	}

	switch {
	case t == "CREATE":
		err = j.mng.Verify(user, settings.CREATE_USER_FIELD_REQUIRED)
	case t == "UPDATE":
		err = j.mng.Verify(user, settings.UPDATE_USER_FIELD_REQUIRED)
	case t == "GET":
		err = j.mng.Verify(user, settings.GET_USER_FIELD_REQUIRED)
	}

	if err != nil {
		return nil, NewApiError(err.Error())
	}

	return user, nil
}
