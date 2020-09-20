package session

import (
	"context"
	"encoding/gob"
	"errors"
	"io"
	"time"
)

// Errors
var (
	// ErrNotFound is the error when session not found
	// store must return ErrNotFound if session data not exists
	ErrNotFound = errors.New("session: not found")
)

// Store interface
type Store interface {
	Get(ctx context.Context, key string) (Data, error)
	Set(ctx context.Context, key string, value Data, opt StoreOption) error
	Del(ctx context.Context, key string) error
}

// StoreOption type
type StoreOption struct {
	TTL time.Duration
}

func makeStoreOption(m *Manager, s *Session) StoreOption {
	return StoreOption{
		TTL: m.config.IdleTimeout,
	}
}

// StoreCoder interface
type StoreCoder interface {
	NewEncoder(w io.Writer) StoreEncoder
	NewDecoder(r io.Reader) StoreDecoder
}

// StoreEncoder interface
type StoreEncoder interface {
	Encode(e interface{}) error
}

// StoreDecoder interface
type StoreDecoder interface {
	Decode(e interface{}) error
}

// DefaultStoreCoder is the default store coder
var DefaultStoreCoder StoreCoder = defaultStoreCoder{}

type defaultStoreCoder struct{}

func (defaultStoreCoder) NewEncoder(w io.Writer) StoreEncoder {
	return gob.NewEncoder(w)
}

func (defaultStoreCoder) NewDecoder(r io.Reader) StoreDecoder {
	return gob.NewDecoder(r)
}
