package redigo_test

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"

	"github.com/acoshift/session"
	store "github.com/acoshift/session/store/redigo"
)

func TestRedis(t *testing.T) {
	t.Parallel()

	s := store.New(store.Config{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}})

	opt := session.StoreOption{TTL: time.Second}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("__redisgo", data, opt)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get("__redisgo", opt)
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set("__redisgo", data, opt)
	time.Sleep(2 * time.Second)
	_, err = s.Get("__redisgo", opt)
	assert.Error(t, err, "expected expired key return error")

	s.Set("__redisgo", data, opt)
	b, err = s.Get("__redisgo", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	_, err = s.Get("__redisgo", session.StoreOption{Rolling: true, TTL: time.Minute})
	assert.NoError(t, err)
	time.Sleep(time.Second)
	_, err = s.Get("__redisgo", opt)
	assert.NoError(t, err)

	s.Del("__redisgo", opt)
	_, err = s.Get("__redisgo", opt)
	assert.Error(t, err)
}

func TestRedisWithoutTTL(t *testing.T) {
	t.Parallel()

	s := store.New(store.Config{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}})

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("__redisgo_without_ttl", data, opt)
	assert.NoError(t, err)

	b, err := s.Get("__redisgo_without_ttl", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
