package store

import (
	"bytes"
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
func (s *Redigo) Get(key string) (session.Data, error) {
	c := s.Pool.Get()
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
func (s *Redigo) Set(key string, value session.Data, opt session.StoreOption) error {
	var buf bytes.Buffer
	err := s.coder().NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}

	c := s.Pool.Get()
	if opt.TTL > 0 {
		_, err = c.Do("SETEX", s.Prefix+key, int64(opt.TTL/time.Second), buf.Bytes())
	} else {
		_, err = c.Do("SET", s.Prefix+key, buf.Bytes())
	}
	c.Close()
	return err
}

// Del deletes session data from redis
func (s *Redigo) Del(key string) error {
	c := s.Pool.Get()
	_, err := c.Do("DEL", s.Prefix+key)
	c.Close()
	return err
}
