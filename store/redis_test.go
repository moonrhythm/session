package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

func redisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return addr
}

func TestRedis(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	s := &Redis{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: redisAddr(),
		}),
	}

	opt := session.StoreOption{TTL: time.Second}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "__redis", data, opt)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get(ctx, "__redis")
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set(ctx, "__redis", data, opt)
	time.Sleep(2 * time.Second)
	_, err = s.Get(ctx, "__redis")
	assert.Error(t, err, "expected expired key return error")

	s.Set(ctx, "__redis", data, opt)
	b, err = s.Get(ctx, "__redis")
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	s.Del(ctx, "__redis")
	_, err = s.Get(ctx, "__redis")
	assert.Error(t, err)
}

func TestRedisWithoutTTL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	s := &Redis{
		Prefix: "session:",
		Client: redis.NewClient(&redis.Options{
			Addr: redisAddr(),
		}),
	}

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "__redis_without_ttl", data, opt)
	assert.NoError(t, err)

	b, err := s.Get(ctx, "__redis_without_ttl")
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
