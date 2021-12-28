package store

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

func TestRedigo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	s := &Redigo{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr())
		},
	}}

	opt := session.StoreOption{TTL: time.Second}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "__redigo", data, opt)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get(ctx, "__redigo")
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set(ctx, "__redigo", data, opt)
	time.Sleep(2 * time.Second)
	_, err = s.Get(ctx, "__redigo")
	assert.Error(t, err, "expected expired key return error")

	s.Set(ctx, "__redigo", data, opt)
	b, err = s.Get(ctx, "__redigo")
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	s.Del(ctx, "__redigo")
	_, err = s.Get(ctx, "__redigo")
	assert.Error(t, err)
}

func TestRedigoWithoutTTL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	s := &Redigo{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr())
		},
	}}

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "__redigo_without_ttl", data, opt)
	assert.NoError(t, err)

	b, err := s.Get(ctx, "__redigo_without_ttl")
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
