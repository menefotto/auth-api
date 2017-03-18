package databases

import (
	"testing"

	"github.com/auth-api/core/databases/amazon"
	"github.com/auth-api/core/databases/google"
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
