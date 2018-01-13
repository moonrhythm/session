package session_test

import (
	"time"

	"github.com/acoshift/session"
)

type mockStore struct {
	GetFunc func(string) (session.Data, error)
	SetFunc func(string, session.Data, time.Duration) error
	DelFunc func(string) error
}

func (m *mockStore) Get(key string) (session.Data, error) {
	if m.GetFunc == nil {
		return nil, nil
	}
	return m.GetFunc(key)
}

func (m *mockStore) Set(key string, value session.Data, ttl time.Duration) error {
	if m.SetFunc == nil {
		return nil
	}
	return m.SetFunc(key, value, ttl)
}

func (m *mockStore) Del(key string) error {
	if m.DelFunc == nil {
		return nil
	}
	return m.DelFunc(key)
}

func (m *mockStore) Touch(key string, ttl time.Duration) error {
	return nil
}
