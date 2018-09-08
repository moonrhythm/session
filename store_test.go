package session_test

import (
	"github.com/moonrhythm/session"
)

type mockStore struct {
	GetFunc func(string, session.StoreOption) (session.Data, error)
	SetFunc func(string, session.Data, session.StoreOption) error
	DelFunc func(string, session.StoreOption) error
}

func (m *mockStore) Get(key string, opt session.StoreOption) (session.Data, error) {
	if m.GetFunc == nil {
		return nil, nil
	}
	return m.GetFunc(key, opt)
}

func (m *mockStore) Set(key string, value session.Data, opt session.StoreOption) error {
	if m.SetFunc == nil {
		return nil
	}
	return m.SetFunc(key, value, opt)
}

func (m *mockStore) Del(key string, opt session.StoreOption) error {
	if m.DelFunc == nil {
		return nil
	}
	return m.DelFunc(key, opt)
}
