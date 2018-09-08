package goredis_test

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
	store "github.com/moonrhythm/session/store/goredis"
)

func TestRedis(t *testing.T) {
	t.Parallel()

	s := store.New(store.Config{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	})

	opt := session.StoreOption{TTL: time.Second}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("__goredis", data, opt)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get("__goredis", opt)
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set("__goredis", data, opt)
	time.Sleep(2 * time.Second)
	_, err = s.Get("__goredis", opt)
	assert.Error(t, err, "expected expired key return error")

	s.Set("__goredis", data, opt)
	b, err = s.Get("__goredis", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	_, err = s.Get("__goredis", session.StoreOption{Rolling: true, TTL: time.Minute})
	assert.NoError(t, err)
	time.Sleep(time.Second)
	_, err = s.Get("__goredis", opt)
	assert.NoError(t, err)

	s.Del("__goredis", opt)
	_, err = s.Get("__goredis", opt)
	assert.Error(t, err)
}

func TestRedisWithoutTTL(t *testing.T) {
	t.Parallel()

	s := store.New(store.Config{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	})

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("__goredis_without_ttl", data, opt)
	assert.NoError(t, err)

	b, err := s.Get("__goredis_without_ttl", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
