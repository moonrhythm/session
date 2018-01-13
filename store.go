package session

import (
	"time"
)

// Store interface
type Store interface {
	Get(key string) (Data, error)
	Set(key string, value Data, ttl time.Duration) error
	Del(key string) error
	Touch(key string, ttl time.Duration) error
}
