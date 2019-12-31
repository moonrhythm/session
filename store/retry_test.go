package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

type mockStore struct {
	attempt int
}

func (s *mockStore) Get(key string, opt session.StoreOption) (session.Data, error) {
	s.attempt++
	if s.attempt == 3 {
		return nil, nil
	}
	return nil, fmt.Errorf("error")
}

func (s *mockStore) Set(key string, value session.Data, opt session.StoreOption) error {
	s.attempt++
	if s.attempt == 3 {
		return nil
	}
	return fmt.Errorf("error")
}

func (s *mockStore) Del(key string, opt session.StoreOption) error {
	s.attempt++
	if s.attempt == 3 {
		return nil
	}
	return fmt.Errorf("error")
}

func TestRetrySuccess(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 0}
		_, err := s.Get("", session.StoreOption{})
		assert.NoError(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 0}
		err := s.Set("", session.Data{}, session.StoreOption{})
		assert.NoError(t, err)
	})

	t.Run("Del", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 0}
		err := s.Del("", session.StoreOption{})
		assert.NoError(t, err)
	})
}

func TestRetryFail(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 1}
		_, err := s.Get("", session.StoreOption{})
		assert.Error(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 1}
		err := s.Set("", session.Data{}, session.StoreOption{})
		assert.Error(t, err)
	})

	t.Run("Del", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 1}
		err := s.Del("", session.StoreOption{})
		assert.Error(t, err)
	})
}
