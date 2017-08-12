package redis_test

import (
	"testing"
	"time"

	store "github.com/acoshift/session/store/redis"
	"github.com/garyburd/redigo/redis"
)

func TestRedis(t *testing.T) {
	s := store.New(store.Config{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}})
	err := s.Set("a", []byte("test"), time.Second)
	if err != nil {
		t.Fatalf("expected set not error; got %v", err)
	}

	time.Sleep(2 * time.Second)
	b, err := s.Get("a")
	if b != nil {
		t.Fatalf("expected expired key return nil value; got %v", b)
	}
	if err == nil {
		t.Fatalf("expected expired key return error; got nil")
	}

	s.Set("a", []byte("test"), time.Second)
	time.Sleep(2 * time.Second)
	_, err = s.Get("a")
	if err == nil {
		t.Fatalf("expected expired key return error; got nil")
	}

	s.Set("a", []byte("test"), time.Second)
	b, err = s.Get("a")
	if err != nil {
		t.Fatalf("expected get valid key return not nil error; got %v", err)
	}
	if string(b) != "test" {
		t.Fatalf("expected get return same value as set")
	}

	s.Del("a")
	_, err = s.Get("a")
	if err == nil {
		t.Fatalf("expected get deleted key to return error; got nil")
	}
}
