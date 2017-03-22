package databases

import (
	"testing"

	"github.com/wind85/auth-api/core/databases/amazon"
	"github.com/wind85/auth-api/core/databases/google"
)

func TestDynamoSatisfyInterface(t *testing.T) {
	var db Db
	db = &amazon.Dynamo{}
	db.Close()
}

func TestDatastoreSatisfyInterface(t *testing.T) {
	var db Db
	db = &google.Datastore{}
	db.Close()
}
