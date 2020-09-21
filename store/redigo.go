package store

import (
	"bytes"
	"context"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/moonrhythm/session"
)

// Redigo is the redis store
// implement by using "github.com/gomodule/redigo/redis"
type Redigo struct {
	Pool   *redis.Pool
	Prefix string
	Coder  session.StoreCoder
}

func (s *Redigo) coder() session.StoreCoder {
	if s.Coder == nil {
		return session.DefaultStoreCoder
	}
	return s.Coder
}

// Get gets session data from redis
func (s *Redigo) Get(ctx context.Context, key string) (session.Data, error) {
	c, err := s.Pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	data, err := redis.Bytes(c.Do("GET", s.Prefix+key))
	c.Close()
	if err == redis.ErrNil {
		return nil, session.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var sessData session.Data
	err = s.coder().NewDecoder(bytes.NewReader(data)).Decode(&sessData)
	if err != nil {
		return nil, err
	}
	return sessData, nil
}

// Set sets session data to redis
func (s *Redigo) Set(ctx context.Context, key string, value session.Data, opt session.StoreOption) error {
	var buf bytes.Buffer
	err := s.coder().NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}

	c, err := s.Pool.GetContext(ctx)
	if err != nil {
		return err
	}
	if opt.TTL > 0 {
		_, err = c.Do("SETEX", s.Prefix+key, int64(opt.TTL/time.Second), buf.Bytes())
	} else {
		_, err = c.Do("SET", s.Prefix+key, buf.Bytes())
	}
	c.Close()
	return err
}

// Del deletes session data from redis
func (s *Redigo) Del(ctx context.Context, key string) error {
	c, err := s.Pool.GetContext(ctx)
	if err != nil {
		return err
	}
	_, err = c.Do("DEL", s.Prefix+key)
	c.Close()
	return err
}
