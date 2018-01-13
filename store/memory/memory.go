package memory

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"
	"time"

	"github.com/acoshift/session"
)

// Config is the memory store config
type Config struct {
	GCInterval time.Duration
}

// New creates new memory store
func New(config Config) session.Store {
	s := &memoryStore{
		gcInterval: config.GCInterval,
		l:          make(map[interface{}]*item),
	}
	if s.gcInterval > 0 {
		time.AfterFunc(s.gcInterval, s.gcWorker)
	}
	return s
}

type item struct {
	data []byte
	exp  time.Time
}

type memoryStore struct {
	gcInterval time.Duration
	m          sync.RWMutex
	l          map[interface{}]*item
}

func (s *memoryStore) gcWorker() {
	s.GC()
	time.AfterFunc(s.gcInterval, s.gcWorker)
}

func (s *memoryStore) GC() {
	now := time.Now()
	s.m.Lock()
	for k, v := range s.l {
		if !v.exp.IsZero() && v.exp.Before(now) {
			delete(s.l, k)
		}
	}
	s.m.Unlock()
}

var errNotFound = errors.New("memory: session not found")

func (s *memoryStore) Get(key string) (session.SessionData, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	v := s.l[key]
	if v == nil {
		return nil, errNotFound
	}
	if !v.exp.IsZero() && v.exp.Before(time.Now()) {
		return nil, errNotFound
	}
	var sessData session.SessionData
	err := gob.NewDecoder(bytes.NewReader(v.data)).Decode(&sessData)
	if err != nil {
		return nil, err
	}
	return sessData, nil
}

func (s *memoryStore) Set(key string, value session.SessionData, ttl time.Duration) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}

	s.m.Lock()
	it := &item{data: buf.Bytes()}
	if ttl > 0 {
		it.exp = time.Now().Add(ttl)
	}
	s.l[key] = it
	s.m.Unlock()
	return nil
}

func (s *memoryStore) Del(key string) error {
	s.m.Lock()
	delete(s.l, key)
	s.m.Unlock()
	return nil
}

func (s *memoryStore) Touch(key string, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}
	s.m.Lock()
	if it, ok := s.l[key]; ok {
		it.exp = time.Now().Add(ttl)
	}
	s.m.Unlock()
	return nil
}
