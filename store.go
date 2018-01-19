package session

import (
	"time"
)

// Store interface
type Store interface {
	Get(key string, opt StoreOption) (Data, error)
	Set(key string, value Data, opt StoreOption) error
	Del(key string, opt StoreOption) error
}

// StoreOption type
type StoreOption struct {
	Rolling bool
	TTL     time.Duration
}
