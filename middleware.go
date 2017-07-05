package session

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/acoshift/middleware"
)

// Middleware is the session parser middleware
func Middleware(config Config) middleware.Middleware {
	if config.Store == nil {
		panic("session: nil store")
	}

	entropy := config.Entropy
	if entropy <= 0 {
		entropy = 16
	}

	name := config.Name
	if len(name) == 0 {
		name = "sess"
	}

	maxAge := int(config.MaxAge / time.Second)

	generateID := func() string {
		b := make([]byte, entropy)
		rand.Read(b)
		return base64.URLEncoding.EncodeToString(b)
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var s Session

			// get session key from cookie
			cookie, err := r.Cookie(name)
			if err == nil && len(cookie.Value) > 0 {
				// get session data from store
				s.p, err = config.Store.Get(cookie.Value)
				if err == nil {
					s.id = cookie.Value
					s.decode(s.p)
				}
			}

			// if session not found, create new session
			if len(s.id) == 0 {
				s.id = generateID()
			}

			// rolling cookie
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Domain:   config.Domain,
				Path:     config.Path,
				HttpOnly: config.HTTPOnly,
				Value:    s.id,
				MaxAge:   maxAge,
				Secure:   (config.Secure == ForceSecure) || (config.Secure == PreferSecure && isTLS(r)),
			})

			// use defer to alway save session even panic
			defer func() {
				if s.markDestory {
					config.Store.Del(s.id)
					return
				}

				// if session was modified, save session to store,
				// if not don't save to store to prevent brute force attack
				b, err := s.encode()
				if err == nil {
					if bytes.Compare(s.p, b) == 0 {
						config.Store.Exp(s.id, config.MaxAge)
						return
					}
					config.Store.Set(s.id, b, config.MaxAge)
				}
			}()

			nr := r.WithContext(Set(r.Context(), &s))
			h.ServeHTTP(w, nr)
		})
	}
}
