package goredis

import (
	"bytes"

	"github.com/go-redis/redis"

	"github.com/moonrhythm/session"
)

// Config is the redis store config
type Config struct {
	Client *redis.Client
	Prefix string
	Coder  session.StoreCoder
}

// New creates new redis store
func New(config Config) session.Store {
	if config.Coder == nil {
		config.Coder = session.DefaultStoreCoder
	}
	return &redisStore{
		client: config.Client,
		prefix: config.Prefix,
		coder:  config.Coder,
	}
}

type redisStore struct {
	client *redis.Client
	prefix string
	coder  session.StoreCoder
}

func (s *redisStore) Get(key string, opt session.StoreOption) (session.Data, error) {
	data, err := s.client.Get(s.prefix + key).Bytes()
	if opt.Rolling && opt.TTL > 0 {
		s.client.Expire(s.prefix+key, opt.TTL)
	}
	if err == redis.Nil {
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
	return s.client.Set(s.prefix+key, buf.Bytes(), opt.TTL).Err()
}

func (s *redisStore) Del(key string, opt session.StoreOption) error {
	return s.client.Del(s.prefix + key).Err()
}
