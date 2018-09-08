package session_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/acoshift/session"
	"github.com/acoshift/session/store/memory"
)

func TestManagerGetSave(t *testing.T) {
	t.Parallel()

	var (
		setKey   string
		setValue session.Data
	)

	m := session.New(session.Config{
		MaxAge: time.Second,
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setKey = key
				setValue = value
				return nil
			},
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				assert.Equal(t, setKey, key)
				return setValue, nil
			},
		},
	})

	h := func(w http.ResponseWriter, r *http.Request) {
		s, _ := m.Get(r, sessName)
		assert.NotEmpty(t, s.ID())
		c, _ := s.Get("test").(int)
		s.Set("test", c+1)
		fmt.Fprintf(w, "%d", c)

		m.Save(w, s)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h(w, r)

	assert.Equal(t, "0", w.Body.String())

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	h(w, r)
	assert.Equal(t, "1", w.Body.String())
}

func TestManagerGetError(t *testing.T) {
	t.Parallel()

	m := session.New(session.Config{
		MaxAge: time.Second,
		Store: &mockStore{
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				return nil, fmt.Errorf("store error")
			},
		},
	})

	h := func(w http.ResponseWriter, r *http.Request) {
		s, err := m.Get(r, sessName)
		assert.Error(t, err)
		assert.Nil(t, s)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sessName+"=test")
	h(w, r)
}

func TestManagerNotPassMiddleware(t *testing.T) {
	t.Parallel()

	m := session.New(session.Config{
		MaxAge: time.Second,
		Store:  memory.New(memory.Config{}),
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	t.Run("Regenerate", func(t *testing.T) {
		s, err := m.Get(r, sessName)
		if assert.NoError(t, err) {
			assert.Error(t, s.Regenerate())
		}
	})

	t.Run("Renew", func(t *testing.T) {
		s, err := m.Get(r, sessName)
		if assert.NoError(t, err) {
			assert.Error(t, s.Renew())
		}
	})

	t.Run("Destroy", func(t *testing.T) {
		s, err := m.Get(r, sessName)
		if assert.NoError(t, err) {
			assert.Error(t, s.Destroy())
		}
	})
}
