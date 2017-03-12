package dynamo

import (
	"fmt"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

const version = 0.1

type Dynamo struct {
	db    *dynamo.DB
	table dynamo.Table
	kind  string
}

func (d *Dynamo) Open(kind, region string) error {
	d.kind = kind
	d.db = dynamo.New(
		session.New(),
		&aws.Config{Region: aws.String(region)},
	)
	if d.db == nil {
		return errors.DynamoInit
	}

	d.table = d.db.Table(kind)

	return nil
}

func (d *Dynamo) Get(key string) (*models.User, error) {
	user := &models.User{}

	err := d.table.Get("Email", key).One(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Dynamo) Put(key string, value *models.User) error {
	_ = key // dumb placeholder
	err := d.table.Put(value).Run()
	if err != nil {
		return err
	}

	return nil
}

func (d *Dynamo) Del(key string) error {
	err := d.table.Delete("Email", key).Run()
	if err != nil {
		return err
	}

	return nil
}

func (d *Dynamo) Backend() *dynamo.DB {
	return d.db
}

func (d *Dynamo) Close() {
	_ = fmt.Sprintf("%s\n", "no op")
}
