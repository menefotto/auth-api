package managers

import (
	"fmt"
	"reflect"
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
	user.Isactive = "false"
	user.Isstaff = "false"
	user.Issuperuser = "false"
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

func (u *Users) Update(newUser *models.User) (*models.User, error) {

	oldUser, err := u.store.Get(newUser.Email)
	if err != nil {
		return nil, err
	}

	if err := u.updateFields(newUser, oldUser); err != nil {
		return nil, err
	}

	return oldUser, nil
}

func (u *Users) Get(user *models.User) (*models.User, error) {
	user, err := u.store.Get(user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Users) updateFields(new, old *models.User) error {

	switch {
	case new.Username != "" && new.Username != old.Username:
		old.Username = new.Username
	case new.Firstname != "" && new.Firstname != old.Firstname:
		old.Firstname = new.Firstname
	case new.Lastname != "" && new.Lastname != old.Lastname:
		old.Lastname = new.Lastname
	case new.Password != "" && new.Password != old.Password:
		old.Password = new.Password
	case new.Email != "" && new.Email != old.Email:
		old.Email = new.Email
	case new.Photourl != "" && new.Photourl != old.Photourl:
		old.Photourl = new.Photourl
	case new.Isactive != "" && new.Isactive != old.Isactive:
		old.Isactive = new.Isactive
	case new.Issuperuser != "" && new.Issuperuser != old.Issuperuser:
		old.Issuperuser = new.Issuperuser
	case new.Isstaff != "" && new.Isstaff != old.Isstaff:
		old.Isstaff = new.Isstaff
	case new.Code != "" && new.Code != old.Code:
		old.Code = new.Code

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
	case fieldname == "Datejoined" || fieldname == "Uuid":
		return nil
	default:
		s, ok := value.(string)
		if !ok {
			return errors.NotString
		}

		if elem.CanSet() {
			elem.SetString(s)
			return nil
		}

	}

	return errors.New("Field: [" + field + "] cannot be set check spelling")
}
