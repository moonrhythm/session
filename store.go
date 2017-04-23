package session

import (
	"time"
)

// Store interface
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, userID interface{}, data []byte, ttl time.Duration) error
	Del(key string) error
	DelUser(userID interface{}) error
}
