package session

import (
	"bytes"
	"context"
	"encoding/gob"
)

// Session type
type Session struct {
	id string
	d  map[interface{}]interface{}
	p  []byte
}

func init() {
	gob.Register(map[interface{}]interface{}{})
}

type sessionKey struct{}

// Get gets session from context
func Get(ctx context.Context) *Session {
	s, _ := ctx.Value(sessionKey{}).(*Session)
	return s
}

// Set sets session to context
func Set(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, s)
}

func (s *Session) encode() ([]byte, error) {
	if len(s.d) == 0 {
		return []byte{}, nil
	}

	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(&s.d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Session) decode(b []byte) {
	s.d = make(map[interface{}]interface{})
	if len(b) > 0 {
		gob.NewDecoder(bytes.NewReader(b)).Decode(&s.d)
	}
}

// Get gets data from session
func (s *Session) Get(key interface{}) interface{} {
	if s.d == nil {
		return nil
	}
	return s.d[key]
}

// Set sets data to session
func (s *Session) Set(key, value interface{}) {
	if s.d == nil {
		s.d = make(map[interface{}]interface{})
	}
	s.d[key] = value
}

// Del deletes data from session
func (s *Session) Del(key interface{}) {
	if s.d == nil {
		return
	}
	delete(s.d, key)
}
