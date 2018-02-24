package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/acoshift/session"
)

func TestSessionRenew(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := session.Get(r.Context(), sessName)
		s.Set("a", 1)
		s.Renew()
		assert.True(t, s.Changed())
		assert.Empty(t, s.GetInt("test"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
}
