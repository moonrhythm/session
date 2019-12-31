package store

import (
	"bytes"

	"github.com/go-redis/redis"

	"github.com/moonrhythm/session"
)

// GoRedis is the redis store
// implement by using "github.com/go-redis/redis" package
type GoRedis struct {
	Client *redis.Client
	Prefix string
	Coder  session.StoreCoder
}

func (s *GoRedis) coder() session.StoreCoder {
	if s.Coder == nil {
		return session.DefaultStoreCoder
	}
	return s.Coder
}

// Get gets session data from redis
func (s *GoRedis) Get(key string, opt session.StoreOption) (session.Data, error) {
	data, err := s.Client.Get(s.Prefix + key).Bytes()
	if opt.Rolling && opt.TTL > 0 {
		s.Client.Expire(s.Prefix+key, opt.TTL)
	}
	if err == redis.Nil {
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
func (s *GoRedis) Set(key string, value session.Data, opt session.StoreOption) error {
	var buf bytes.Buffer
	err := s.coder().NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}
	return s.Client.Set(s.Prefix+key, buf.Bytes(), opt.TTL).Err()
}

// Del deletes session data from redis
func (s *GoRedis) Del(key string, opt session.StoreOption) error {
	return s.Client.Del(s.Prefix + key).Err()
}
