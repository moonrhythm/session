package redigo

import (
	"bytes"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/moonrhythm/session"
)

// Config is the redis store config
type Config struct {
	Pool   *redis.Pool
	Prefix string
	Coder  session.StoreCoder
}

// New creates new redis store
func New(config Config) session.Store {
	if config.Coder == nil {
		config.Coder = session.DefaultStoreCoder
	}
	return &redisStore{
		pool:   config.Pool,
		prefix: config.Prefix,
		coder:  config.Coder,
	}
}

type redisStore struct {
	pool   *redis.Pool
	prefix string
	coder  session.StoreCoder
}

func (s *redisStore) Get(key string, opt session.StoreOption) (session.Data, error) {
	c := s.pool.Get()
	data, err := redis.Bytes(c.Do("GET", s.prefix+key))
	if opt.Rolling && opt.TTL > 0 {
		c.Do("EXPIRE", s.prefix+key, int64(opt.TTL/time.Second))
	}
	c.Close()
	if err == redis.ErrNil {
		return nil, session.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var sessData session.Data
	err = s.coder.NewDecoder(bytes.NewReader(data)).Decode(&sessData)
	if err != nil {
		return nil, err
	}
	return sessData, nil
}

func (s *redisStore) Set(key string, value session.Data, opt session.StoreOption) error {
	var buf bytes.Buffer
	err := s.coder.NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}

	c := s.pool.Get()
	if opt.TTL > 0 {
		_, err = c.Do("SETEX", s.prefix+key, int64(opt.TTL/time.Second), buf.Bytes())
	} else {
		_, err = c.Do("SET", s.prefix+key, buf.Bytes())
	}
	c.Close()
	return err
}

func (s *redisStore) Del(key string, opt session.StoreOption) error {
	c := s.pool.Get()
	_, err := c.Do("DEL", s.prefix+key)
	c.Close()
	return err
}
