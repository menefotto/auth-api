package managers

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/auth-api/core/google"
	"github.com/auth-api/core/models"
	"github.com/pborman/uuid"

	"github.com/auth-api/core/settings"
)

var ErrNotBool = errors.New("Not a bool value")
var ErrNotString = errors.New("Not a string value")

type Users struct {
	store *google.Datastore
}

func New(kind string) *Users {
	db := &google.Datastore{}

	err := db.Open(settings.PROJECTID, kind)
	if err != nil {
		panic(err)
	}

	return &Users{db}
}

func (u *Users) Create(user *models.User) (*models.User, error) {

	user.Uuid = uuid.New()
	user.IsActive = false
	user.IsStaff = false
	user.IsSuperUser = false
	user.DateJoined = fmt.Sprint(time.Now().UTC())

	if err := u.store.Put(user.Email, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Users) Update(fields map[string]interface{}) (*models.User, error) {

	email, ok := fields["Email"].(string)
	if !ok {
		return nil, ErrNotString
	}

	oldUser, err := u.store.Get(email)
	if err != nil {
		return nil, err
	}

	if err := u.updateFields(fields, oldUser); err != nil {
		return nil, err
	}

	return oldUser, nil
}

func (u *Users) Get(email string) (*models.User, error) {
	user, err := u.store.Get(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Users) updateFields(mapped map[string]interface{}, old *models.User) error {

	for k, value := range mapped {
		if err := SetField(old, k, value); err != nil {
			return err
		}
	}

	if err := u.store.Put(old.Email, old); err != nil {
		return err
	}

	return nil
}

func (u *Users) Verify(user *models.User, settings []string) error {

	for _, field := range settings {
		value := getField(user, field)
		if len(value) < 2 {
			return errors.New(field + " is required!")
		}
	}

	return nil
}

func getField(user *models.User, field string) string {
	value := reflect.ValueOf(user).Elem().FieldByName(field)

	return fmt.Sprint(value)
}

func SetField(user *models.User, field string, value interface{}) error {

	elem := reflect.ValueOf(user).Elem().FieldByName(field)

	switch {
	case field == "isactive" || field == "issuperuser" || field == "isstaff":
		b, ok := value.(bool)
		if !ok {
			return ErrNotBool
		}

		elem.SetBool(b)
	default:
		s, ok := value.(string)
		if !ok {
			return ErrNotString
		}

		elem.SetString(s)
	}

	return nil
}
