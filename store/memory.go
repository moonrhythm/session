package store

import (
	"bytes"
	"sync"
	"time"

	"github.com/moonrhythm/session"
)

// Memory stores session data in memory
type Memory struct {
	Coder session.StoreCoder

	m sync.RWMutex
	l map[interface{}]*memoryItem
}

type memoryItem struct {
	data []byte
	exp  time.Time
}

func (s *Memory) coder() session.StoreCoder {
	if s.Coder == nil {
		return session.DefaultStoreCoder
	}
	return s.Coder
}

func (s *Memory) gcWorker(d time.Duration) {
	s.GC()
	time.AfterFunc(d, func() { s.gcWorker(d) })
}

// GCEvery starts gc every given duration
func (s *Memory) GCEvery(d time.Duration) *Memory {
	time.AfterFunc(d, func() { s.gcWorker(d) })
	return s
}

// GC runs gc
func (s *Memory) GC() {
	s.m.Lock()
	defer s.m.Unlock()

	now := time.Now()
	for k, v := range s.l {
		if !v.exp.IsZero() && v.exp.Before(now) {
			delete(s.l, k)
		}
	}
}

// Get gets session data from memory
func (s *Memory) Get(key string, opt session.StoreOption) (session.Data, error) {
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
	err := s.coder().NewDecoder(bytes.NewReader(v.data)).Decode(&sessData)
	if err != nil {
		return nil, err
	}
	if opt.Rolling && opt.TTL > 0 {
		v.exp = time.Now().Add(opt.TTL)
	}
	return sessData, nil
}

// Set sets session data to memory
func (s *Memory) Set(key string, value session.Data, opt session.StoreOption) error {
	var buf bytes.Buffer
	err := s.coder().NewEncoder(&buf).Encode(value)
	if err != nil {
		return err
	}

	s.m.Lock()
	it := &memoryItem{data: buf.Bytes()}
	if opt.TTL > 0 {
		it.exp = time.Now().Add(opt.TTL)
	}
	if s.l == nil {
		s.l = make(map[interface{}]*memoryItem)
	}
	s.l[key] = it
	s.m.Unlock()
	return nil
}

// Del deletes session data from memory
func (s *Memory) Del(key string, opt session.StoreOption) error {
	s.m.Lock()
	delete(s.l, key)
	s.m.Unlock()
	return nil
}
