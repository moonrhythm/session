package session

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

// Manager is the session manager
type Manager struct {
	config Config
	hashID func(id string) string
}

// manager internal data
const (
	timestampKey = "_session/timestamp"
	destroyedKey = "_session/destroyed" // for detect session hijack
)

// New creates new session manager
func New(config Config) *Manager {
	if config.Store == nil {
		panic("session: nil store")
	}

	m := Manager{}
	m.config = config

	m.hashID = func(id string) string {
		h := sha256.New()
		h.Write([]byte(id))
		h.Write(config.Secret)
		return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")
	}

	if config.DisableHashID {
		m.hashID = func(id string) string {
			return id
		}
	}

	return &m
}

// Get retrieves session from request
func (m *Manager) Get(r *http.Request, name string) *Session {
	s := Session{
		Name:       name,
		Domain:     m.config.Domain,
		Path:       m.config.Path,
		HTTPOnly:   m.config.HTTPOnly,
		MaxAge:     m.config.MaxAge,
		Secure:     (m.config.Secure == ForceSecure) || (m.config.Secure == PreferSecure && isTLS(r, m.config.TrustProxy)),
		SameSite:   m.config.SameSite,
		Rolling:    m.config.Rolling,
		IDHashFunc: m.hashID,
	}

	// get session key from cookie
	cookie, err := r.Cookie(name)
	if err == nil && len(cookie.Value) > 0 {
		hashedID := m.hashID(cookie.Value)

		// get session data from store
		s.data, err = m.config.Store.Get(hashedID, makeStoreOption(m, &s))
		if err == nil {
			s.rawID = cookie.Value
			s.id = hashedID
		}
		// DO NOT set session id to cookie value if not found in store
		// to prevent session fixation attack
	}

	if len(s.id) == 0 {
		s.rawID = generateID()
		s.id = m.hashID(s.rawID)
		s.isNew = true
	}

	return &s
}

// Save saves session to store and set cookie to response
//
// Save must be called before response header was written
func (m *Manager) Save(w http.ResponseWriter, s *Session) error {
	s.setCookie(w)

	opt := makeStoreOption(m, s)

	if s.destroy {
		m.config.Store.Del(s.id, opt)
		return nil
	}

	// detect is flash changed and encode new flash data
	if s.flash != nil && s.flash.Changed() {
		b, _ := s.flash.Encode()
		s.Set(flashKey, b)
	}

	// if session not modified, don't save to store to prevent store overflow
	if !s.Changed() {
		return nil
	}

	// check is regenerate
	if len(s.oldID) > 0 {
		if m.config.DeleteOldSession {
			m.config.Store.Del(s.oldID, opt)
		} else {
			// save old session data if not delete
			s.oldData[timestampKey] = int64(0)
			s.oldData[destroyedKey] = time.Now().UnixNano()
			err := m.config.Store.Set(s.oldID, s.oldData, opt)
			if err != nil {
				return err
			}
		}
	}

	// save sesion data to store
	s.Set(timestampKey, time.Now().Unix())
	err := m.config.Store.Set(s.id, s.data, opt)
	return err
}
