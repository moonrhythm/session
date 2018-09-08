package session

import (
	"context"
	"errors"
	"net/http"
)

// Errors
var (
	ErrNotPassMiddleware = errors.New("session: request not pass middleware")
)

type managerKey struct{}

// Middleware is the Manager middleware wrapper
//
// New(config).Middleware()
func Middleware(config Config) func(http.Handler) http.Handler {
	return New(config).Middleware()
}

// Middleware injects session manager into request's context.
//
// All data changed before write response writer's header will be save.
func (m *Manager) Middleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rm := requestedManager{
				Manager: m,
				r:       r,
				storage: make(map[string]*Session),
			}

			ctx := context.WithValue(r.Context(), managerKey{}, &rm)
			nr := r.WithContext(ctx)
			nw := sessionWriter{
				ResponseWriter: w,
				beforeWriteHeader: func() {
					for _, s := range rm.storage {
						err := m.Save(w, s)
						if err != nil {
							panic("session: " + err.Error())
						}
					}
				},
			}
			h.ServeHTTP(&nw, nr)

			if !nw.wroteHeader {
				nw.beforeWriteHeader()
			}
		})
	}
}

// Get gets session from context
func Get(ctx context.Context, name string) (*Session, error) {
	m, _ := ctx.Value(managerKey{}).(*requestedManager)
	if m == nil {
		return nil, ErrNotPassMiddleware
	}

	// try get session from storage first
	// to preserve session data from difference handler
	if s, ok := m.storage[name]; ok {
		return s, nil
	}

	// get session from manager
	s, err := m.Get(name)
	if err != nil {
		return nil, err
	}

	// save session to storage for later get
	m.storage[name] = s
	return s, nil
}

type requestedManager struct {
	*Manager
	r       *http.Request
	storage map[string]*Session
}

func (m *requestedManager) Get(name string) (*Session, error) {
	return m.Manager.Get(m.r, name)
}
