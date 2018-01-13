package redis_test

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"

	"github.com/acoshift/session"
	store "github.com/acoshift/session/store/redis"
)

func TestRedis(t *testing.T) {
	s := store.New(store.Config{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}})

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("a", data, time.Second)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get("a")
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set("a", data, time.Second)
	time.Sleep(2 * time.Second)
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

func TestRedisWithoutMaxAge(t *testing.T) {
	s := store.New(store.Config{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}})

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("a", data, 0)
	assert.NoError(t, err)

	b, err := s.Get("a")
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
