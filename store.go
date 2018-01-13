package session

import (
	"time"
)

// Store interface
type Store interface {
	Get(key string) (SessionData, error)
	Set(key string, value SessionData, ttl time.Duration) error
	Del(key string) error
	Touch(key string, ttl time.Duration) error
}
