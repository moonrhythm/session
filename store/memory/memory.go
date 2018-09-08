package memory

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/moonrhythm/session"
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

func (s *memoryStore) Get(key string, opt session.StoreOption) (session.Data, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	v := s.l[key]
	if v == nil {
		return nil, session.ErrNotFound
	}
	if !v.exp.IsZero() && v.exp.Before(time.Now()) {
		return nil, session.ErrNotFound
	}
	var sessData session.Data
	err := gob.NewDecoder(bytes.NewReader(v.data)).Decode(&sessData)
	if err != nil {
		return nil, err
	}
	if opt.Rolling && opt.TTL > 0 {
		v.exp = time.Now().Add(opt.TTL)
	}
	return sessData, nil
}

func (s *memoryStore) Set(key string, value session.Data, opt session.StoreOption) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}

	s.m.Lock()
	it := &item{data: buf.Bytes()}
	if opt.TTL > 0 {
		it.exp = time.Now().Add(opt.TTL)
	}
	s.l[key] = it
	s.m.Unlock()
	return nil
}

func (s *memoryStore) Del(key string, opt session.StoreOption) error {
	s.m.Lock()
	delete(s.l, key)
	s.m.Unlock()
	return nil
}
