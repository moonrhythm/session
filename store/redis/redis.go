package redis

import (
	"time"

	"github.com/acoshift/session"
	"github.com/garyburd/redigo/redis"
)

// New creates new redis store
func New(pool *redis.Pool) session.Store {
	return &redisStore{pool}
}

type redisStore struct {
	pool *redis.Pool
}

func (s *redisStore) Get(key string) ([]byte, error) {
	c := s.pool.Get()
	defer c.Close()
	return redis.Bytes(c.Do("GET", key))
}

func (s *redisStore) Set(key string, value []byte, ttl time.Duration) error {
	c := s.pool.Get()
	defer c.Close()
	_, err := c.Do("SETEX", key, int64(ttl/time.Second), data)
	return err
}

func (s *redisStore) Del(key string) error {
	c := s.pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", key)
	return err
}

func (s *redisStore) Exp(key string, ttl time.Duration) error {
	c := s.pool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", key, int64(ttl/time.Second))
	return err
}
