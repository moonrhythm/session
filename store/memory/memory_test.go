package memory_test

import (
	"testing"
	"time"

	"github.com/acoshift/session"
	store "github.com/acoshift/session/store/memory"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	s := store.New(store.Config{GCInterval: 10 * time.Millisecond})

	data := make(session.SessionData)
	data["test"] = "123"

	err := s.Set("a", data, time.Millisecond)
	assert.NoError(t, err)

	time.Sleep(5 * time.Millisecond)
	b, err := s.Get("a")
	assert.Nil(t, b)
	assert.Error(t, err)

	s.Set("a", data, time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, err = s.Get("a")
	assert.Error(t, err, "expected expired key return error")

	s.Set("a", data, time.Second)
	b, err = s.Get("a")
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	err = s.Touch("a", time.Minute)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	_, err = s.Get("a")
	assert.NoError(t, err)

	s.Del("a")
	_, err = s.Get("a")
	assert.Error(t, err)
}

func TestMemoryWithoutTTL(t *testing.T) {
	s := store.New(store.Config{GCInterval: 10 * time.Millisecond})

	data := make(session.SessionData)
	data["test"] = "123"

	err := s.Set("a", data, 0)
	assert.NoError(t, err)

	b, err := s.Get("a")
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
