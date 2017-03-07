package managers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/google"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/utils"
	"github.com/pborman/uuid"

	"github.com/auth-api/core/settings"
)

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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, err
	}

	user.Password = string(hash)
	user.Uuid = uuid.New()
	user.Isactive = false
	user.Isstaff = false
	user.Issuperuser = false
	user.Datejoined = fmt.Sprint(time.Now().UTC())
	user.Code = utils.GenerateJwt(
		[]byte(user.Email),
		settings.JWT_ACTIVATION_DELTA,
	)

	if err := u.store.Put(user.Email, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Users) Update(fields map[string]interface{}) (*models.User, error) {

	email, ok := fields["email"].(string)
	if !ok {
		return nil, errors.NotString
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

	fieldname := strings.ToUpper(field[:1]) + field[1:]
	elem := reflect.ValueOf(user).Elem().FieldByName(fieldname)

	switch {
	case fieldname == "Isactive" || fieldname == "Issuperuser" || fieldname == "Isstaff":
		s, ok := value.(string)
		if !ok {
			return errors.NotBool
		}

		b, err := strconv.ParseBool(s)
		if err != nil {
			return errors.NotBool
		}

		if elem.CanSet() {
			elem.SetBool(b)
			return nil
		}

	case fieldname == "Email" || fieldname == "Photourl" ||
		fieldname == "Firstname" || fieldname == "Lastname" ||
		fieldname == "Password" || fieldname == "Username" ||
		fieldname == "Code" || fieldname == "Datajoined":

		s, ok := value.(string)
		if !ok {
			return errors.NotString
		}

		if elem.CanSet() {
			elem.SetString(s)
			return nil
		}

	}

	return errors.New("Field: [" + strings.ToLower(fieldname) + "] cannot be set check spelling")
}
