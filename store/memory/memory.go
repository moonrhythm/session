package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/acoshift/session"
)

// New creates new memory store
func New() session.Store {
	s := &memoryStore{
		l: make(map[interface{}]*item),
	}
	go s.cleanupWorker()
	return s
}

type item struct {
	data []byte
	exp  time.Time
}

type memoryStore struct {
	m sync.RWMutex
	l map[interface{}]*item
}

func (s *memoryStore) cleanupWorker() {
	time.Sleep(5 * time.Second)
	for {
		now := time.Now()
		s.m.Lock()
		for k, v := range s.l {
			if !v.exp.IsZero() && v.exp.Before(now) {
				delete(s.l, k)
			}
		}
		s.m.Unlock()
		time.Sleep(6 * time.Hour)
	}
}

var errNotFound = errors.New("memory: session not found")

func (s *memoryStore) Get(key string) ([]byte, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	v := s.l[key]
	if v == nil {
		return nil, errNotFound
	}
	if !v.exp.IsZero() && v.exp.Before(time.Now()) {
		return nil, errNotFound
	}
	return v.data, nil
}

func (s *memoryStore) Set(key string, value []byte, ttl time.Duration) error {
	s.m.Lock()
	s.l[key] = &item{
		data: value,
		exp:  time.Now().Add(ttl),
	}
	s.m.Unlock()
	return nil
}

func (s *memoryStore) Del(key string) error {
	s.m.Lock()
	delete(s.l, key)
	s.m.Unlock()
	return nil
}

func (s *memoryStore) Exp(key string, ttl time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()
	v := s.l[key]
	if v == nil {
		return nil
	}
	v.exp = time.Now().Add(ttl)
	return nil
}