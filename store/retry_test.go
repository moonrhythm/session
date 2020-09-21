package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

type mockStore struct {
	attempt int
}

func (s *mockStore) Get(ctx context.Context, key string) (session.Data, error) {
	s.attempt++
	if s.attempt == 3 {
		return nil, nil
	}
	return nil, fmt.Errorf("error")
}

func (s *mockStore) Set(ctx context.Context, key string, value session.Data, opt session.StoreOption) error {
	s.attempt++
	if s.attempt == 3 {
		return nil
	}
	return fmt.Errorf("error")
}

func (s *mockStore) Del(ctx context.Context, key string) error {
	s.attempt++
	if s.attempt == 3 {
		return nil
	}
	return fmt.Errorf("error")
}

func TestRetrySuccess(t *testing.T) {
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 0}
		_, err := s.Get(ctx, "")
		assert.NoError(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 0}
		err := s.Set(ctx, "", session.Data{}, session.StoreOption{})
		assert.NoError(t, err)
	})

	t.Run("Del", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 0}
		err := s.Del(ctx, "")
		assert.NoError(t, err)
	})
}

func TestRetryFail(t *testing.T) {
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 1}
		_, err := s.Get(ctx, "")
		assert.Error(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 1}
		err := s.Set(ctx, "", session.Data{}, session.StoreOption{})
		assert.Error(t, err)
	})

	t.Run("Del", func(t *testing.T) {
		s := &Retry{Store: &mockStore{}, MaxAttempts: 1}
		err := s.Del(ctx, "")
		assert.Error(t, err)
	})
}
