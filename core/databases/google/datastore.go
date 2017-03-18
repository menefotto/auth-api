package google

import (
	"errors"

	"golang.org/x/net/context"

	"cloud.google.com/go/datastore"
	"github.com/auth-api/core/models"
)

const version = 0.1

//Datastore is the what contains the underlining db handle
type Datastore struct {
	db   *datastore.Client
	kind string
}

// errors
var (
	ErrBucketNotFound   = errors.New("Bucket not found")
	ErrBucketCantCreate = errors.New("Could not create bucket")
	ErrKeyNotFound      = errors.New("Key not found")
	ErrDoesNotExist     = errors.New("Does not exist")
	ErrFoundIt          = errors.New("Found it")
	ErrExistsInSet      = errors.New("Element already exists in set")
	ErrInvalidID        = errors.New("Element ID can not contain \":\"")
	ErrCantDelete       = errors.New("Could not delete key")
)

// Open the datastore
func (d *Datastore) Open(projectID, kind string) error {
	var err error

	d.kind = kind
	d.db, err = datastore.NewClient(context.Background(), projectID)
	if err != nil {
		return err
	}

	return nil
}

// Get the keys as the name suggest
func (d *Datastore) Get(key string) (*models.User, error) {
	user := &models.User{}

	err := d.db.Get(context.Background(), d.NewKey(key), user)
	if err != nil {
		return nil, ErrKeyNotFound
	}

	return user, nil
}

// Put keys inside as the name suggest
func (d *Datastore) Put(key string, value *models.User) error {
	_, err := d.db.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
			if _, err := tx.Put(d.NewKey(key), value); err != nil {
				return err
			}

			return nil
		})

	// maybe do something here

	return err
}

// Del delete keys
func (d *Datastore) Del(key string) error {
	err := d.db.Delete(context.Background(), d.NewKey(key))
	if err != nil {
		return ErrCantDelete
	}

	return nil
}

// Backend give back the db backend implementation
func (d *Datastore) Backend() *datastore.Client {
	return d.db
}

// Close always returns
func (d *Datastore) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d *Datastore) NewKey(id string) *datastore.Key {
	return datastore.NameKey(d.kind, id, nil)
}
