package store

import (
	"bytes"
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/moonrhythm/session"
)

// Redis is the redis store
// implement by using "github.com/redis/go-redis/v9" package
type Redis struct {
	Client *redis.Client
	Prefix string
	Coder  session.StoreCoder
}

func (s *Redis) coder() session.StoreCoder {
	if s.Coder == nil {
		return session.DefaultStoreCoder
	}
	return s.Coder
}

// Get gets session data from redis
func (s *Redis) Get(ctx context.Context, key string) (session.Data, error) {
	data, err := s.Client.Get(ctx, s.Prefix+key).Bytes()
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
func (s *Redis) Set(ctx context.Context, key string, value session.Data, opt session.StoreOption) error {
	var buf bytes.Buffer
	err := s.coder().NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, s.Prefix+key, buf.Bytes(), opt.TTL).Err()
}

// Del deletes session data from redis
func (s *Redis) Del(ctx context.Context, key string) error {
	return s.Client.Del(ctx, s.Prefix+key).Err()
}
