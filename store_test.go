package session_test

import (
	"github.com/moonrhythm/session"
)

type mockStore struct {
	GetFunc func(string) (session.Data, error)
	SetFunc func(string, session.Data, session.StoreOption) error
	DelFunc func(string) error
}

func (m *mockStore) Get(key string) (session.Data, error) {
	if m.GetFunc == nil {
		return nil, nil
	}
	return m.GetFunc(key)
}

func (m *mockStore) Set(key string, value session.Data, opt session.StoreOption) error {
	if m.SetFunc == nil {
		return nil
	}
	return m.SetFunc(key, value, opt)
}

func (m *mockStore) Del(key string) error {
	if m.DelFunc == nil {
		return nil
	}
	return m.DelFunc(key)
}
