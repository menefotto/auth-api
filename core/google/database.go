package google

import "github.com/api.jwt.auth/core/models"

// Db interface implemented by clients
type Db interface {
	Open(projectID, kind string) error
	Get(key string) (*models.User, error)
	Put(key string, value *models.User) error
	Del(key string) error
	Close()
}
