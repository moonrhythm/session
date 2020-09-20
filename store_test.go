package session_test

import (
	"context"

	"github.com/moonrhythm/session"
)

type mockStore struct {
	GetFunc func(string) (session.Data, error)
	SetFunc func(string, session.Data, session.StoreOption) error
	DelFunc func(string) error
}

func (m *mockStore) Get(ctx context.Context, key string) (session.Data, error) {
	if m.GetFunc == nil {
		return nil, nil
	}
	return m.GetFunc(key)
}

func (m *mockStore) Set(ctx context.Context, key string, value session.Data, opt session.StoreOption) error {
	if m.SetFunc == nil {
		return nil
	}
	return m.SetFunc(key, value, opt)
}

func (m *mockStore) Del(ctx context.Context, key string) error {
	if m.DelFunc == nil {
		return nil
	}
	return m.DelFunc(key)
}
