package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

func TestMemory(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	s := new(Memory).GCEvery(10 * time.Millisecond)

	opt := session.StoreOption{TTL: time.Millisecond}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "a", data, opt)
	assert.NoError(t, err)

	time.Sleep(5 * time.Millisecond)
	b, err := s.Get(ctx, "a")
	assert.Nil(t, b)
	assert.Error(t, err)

	s.Set(ctx, "a", data, opt)
	time.Sleep(20 * time.Millisecond)
	_, err = s.Get(ctx, "a")
	assert.Error(t, err, "expected expired key return error")

	s.Set(ctx, "a", data, session.StoreOption{TTL: time.Second})
	b, err = s.Get(ctx, "a")
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	s.Del(ctx, "a")
	_, err = s.Get(ctx, "a")
	assert.Error(t, err)
}

func TestMemoryWithoutTTL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	s := new(Memory).GCEvery(10 * time.Millisecond)

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "a", data, opt)
	assert.NoError(t, err)

	b, err := s.Get(ctx, "a")
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
