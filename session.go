package session

import (
	"net/http"
	"time"

	"github.com/acoshift/flash"
)

// Data stores session data
type Data map[string]interface{}

// Session type
type Session struct {
	id      string // id is the hashed id if enable hash
	rawID   string
	oldID   string // for regenerate, is the hashed old id if enable hash
	oldData Data   // is the old encoded data before regenerate
	data    Data
	destroy bool
	changed bool
	flash   *flash.Flash

	// cookie config
	Name     string
	Domain   string
	Path     string
	HTTPOnly bool
	MaxAge   time.Duration
	Secure   bool
	SameSite SameSite

	IDHashFunc func(id string) string
}

// Clone clones session data
func (data Data) Clone() Data {
	r := make(Data)
	for k, v := range data {
		r[k] = v
	}
	return r
}

// session internal data
const (
	flashKey = "session/flash"
)

// ID returns session id or hashed session id if enable hash id
func (s *Session) ID() string {
	return s.id
}

// Changed returns is session data changed
func (s *Session) Changed() bool {
	if s.changed {
		return true
	}
	if s.flash != nil && s.flash.Changed() {
		s.changed = true
		return true
	}
	return false
}

// Get gets data from session
func (s *Session) Get(key string) interface{} {
	if s.data == nil {
		return nil
	}
	return s.data[key]
}

// Set sets data to session
func (s *Session) Set(key string, value interface{}) {
	if s.data == nil {
		s.data = make(Data)
	}
	s.changed = true
	s.data[key] = value
}

// Del deletes data from session
func (s *Session) Del(key string) {
	if s.data == nil {
		return
	}
	if _, ok := s.data[key]; ok {
		s.changed = true
		delete(s.data, key)
	}
}

// Pop gets data from session then delete it
func (s *Session) Pop(key string) interface{} {
	if s.data == nil {
		return nil
	}
	r := s.data[key]
	s.changed = true
	delete(s.data, key)
	return r
}

// Regenerate regenerates session id
// use when change user access level to prevent session fixation
//
// can not use regenerate and destroy same time
// Regenerate can call only one time
func (s *Session) Regenerate() {
	if len(s.oldID) > 0 {
		return
	}

	if s.destroy {
		return
	}

	s.oldID = s.id
	s.oldData = s.data.Clone()
	s.rawID = generateID()
	if s.IDHashFunc != nil {
		s.id = s.IDHashFunc(s.rawID)
	} else {
		s.id = s.rawID
	}
	s.changed = true
}

// Renew clear all data in current session
// and regenerate session id
func (s *Session) Renew() {
	s.data = make(Data)
	s.Regenerate()
}

// Destroy destroys session from store
func (s *Session) Destroy() {
	s.destroy = true
}

func (s *Session) setCookie(w http.ResponseWriter) {
	if s.destroy {
		http.SetCookie(w, &http.Cookie{
			Name:     s.Name,
			Domain:   s.Domain,
			Path:     s.Path,
			HttpOnly: s.HTTPOnly,
			Value:    "",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			Secure:   s.Secure,
		})
		return
	}

	// if session don't have raw id, don't set cookie
	if len(s.rawID) == 0 {
		return
	}

	// if session not modified, don't set cookie
	if !s.Changed() {
		return
	}

	setCookie(w, &cookie{
		Cookie: http.Cookie{
			Name:     s.Name,
			Domain:   s.Domain,
			Path:     s.Path,
			HttpOnly: s.HTTPOnly,
			Value:    s.rawID,
			MaxAge:   int(s.MaxAge / time.Second),
			Expires:  time.Now().Add(s.MaxAge),
			Secure:   s.Secure,
		},
		SameSite: s.SameSite,
	})
}

// Flash returns flash from session
func (s *Session) Flash() *flash.Flash {
	if s.flash != nil {
		return s.flash
	}
	if b, ok := s.Get(flashKey).([]byte); ok {
		s.flash, _ = flash.Decode(b)
	}
	if s.flash == nil {
		s.flash = flash.New()
	}
	return s.flash
}

// Hijacked checks is session was hijacked,
// can use only with Manager
func (s *Session) Hijacked() bool {
	if t, ok := s.Get(destroyedKey).(int64); ok {
		if t < time.Now().UnixNano()-int64(HijackedTime) {
			return true
		}
	}
	return false
}
