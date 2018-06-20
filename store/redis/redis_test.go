package redis_test

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"

	"github.com/acoshift/session"
	store "github.com/acoshift/session/store/redis"
)

func TestRedis(t *testing.T) {
	s := store.New(store.Config{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	})

	opt := session.StoreOption{TTL: time.Second}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("a", data, opt)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get("a", opt)
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set("a", data, opt)
	time.Sleep(2 * time.Second)
	_, err = s.Get("a", opt)
	assert.Error(t, err, "expected expired key return error")

	s.Set("a", data, opt)
	b, err = s.Get("a", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	_, err = s.Get("a", session.StoreOption{Rolling: true, TTL: time.Minute})
	assert.NoError(t, err)
	time.Sleep(time.Second)
	_, err = s.Get("a", opt)
	assert.NoError(t, err)

	s.Del("a", opt)
	_, err = s.Get("a", opt)
	assert.Error(t, err)
}

func TestRedisWithoutTTL(t *testing.T) {
	s := store.New(store.Config{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	})

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("a", data, opt)
	assert.NoError(t, err)

	b, err := s.Get("a", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
