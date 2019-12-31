package store

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

func TestGoRedis(t *testing.T) {
	t.Parallel()

	s := &GoRedis{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}

	opt := session.StoreOption{TTL: time.Second}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("__goredis", data, opt)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get("__goredis")
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set("__goredis", data, opt)
	time.Sleep(2 * time.Second)
	_, err = s.Get("__goredis")
	assert.Error(t, err, "expected expired key return error")

	s.Set("__goredis", data, opt)
	b, err = s.Get("__goredis")
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	s.Del("__goredis")
	_, err = s.Get("__goredis")
	assert.Error(t, err)
}

func TestGoRedisWithoutTTL(t *testing.T) {
	t.Parallel()

	s := &GoRedis{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set("__goredis_without_ttl", data, opt)
	assert.NoError(t, err)

	b, err := s.Get("__goredis_without_ttl")
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
