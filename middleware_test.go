package session_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/acoshift/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/acoshift/session"
	"github.com/acoshift/session/store/memory"
)

const sessName = "sess"

func mockHandlerFunc(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context(), sessName)
	s.Set("test", 1)
	w.Write([]byte("ok"))
}

var mockHandler = http.HandlerFunc(mockHandlerFunc)

func TestPanicConfig(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err, "expected panic when misconfig")
	}()
	session.Middleware(session.Config{})
}

func TestDefaultConfig(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{},
	})(mockHandler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	cookie := w.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie, "expected cookie not empty")
}

func TestEmptySession(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				assert.Fail(t, "expected get was not called")
				return nil, nil
			},
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				assert.Fail(t, "expected set was not called")
				return nil
			},
			DelFunc: func(key string, opt session.StoreOption) error {
				assert.Fail(t, "expected del was not called")
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	cookie := w.Header().Get("Set-Cookie")
	assert.Empty(t, cookie, "expected cookie empty")
}

func TestEmptySessionFlash(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				assert.Fail(t, "expected get was not called")
				return nil, nil
			},
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				assert.Fail(t, "expected set was not called")
				return nil
			},
			DelFunc: func(key string, opt session.StoreOption) error {
				assert.Fail(t, "expected del was not called")
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		f := s.Flash()
		f.Clear()
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	cookie := w.Header().Get("Set-Cookie")
	assert.Empty(t, cookie, "expected cookie empty")
}

func TestSessionSetInStore(t *testing.T) {
	var (
		setCalled bool
		setKey    string
		setValue  session.Data
		setTTL    time.Duration
	)

	h := session.Middleware(session.Config{
		MaxAge: time.Second,
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setCalled = true
				setKey = key
				setValue = value
				setTTL = opt.TTL
				return nil
			},
		},
	})(mockHandler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	assert.True(t, setCalled, "expected store was called")
	assert.NotEmpty(t, setKey, "expected key not empty")
	assert.NotEmpty(t, setValue, "expected value not empty")
	assert.Equal(t, time.Second, setTTL)

	cs := w.Result().Cookies()
	assert.Len(t, cs, 1, "expected response has 1 cookie; got %d", len(cs))
	assert.NotEqual(t, setKey, cs[0].Value, "expected session id was hashed")
}

func TestSessionGetSet(t *testing.T) {
	var (
		setCalled int
		setKey    string
		setValue  session.Data
	)

	h := session.Middleware(session.Config{
		MaxAge: time.Second,
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setCalled++
				setKey = key
				setValue = value
				return nil
			},
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				assert.Equal(t, setKey, key)
				return setValue, nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		assert.NotEmpty(t, s.ID())
		c, _ := s.Get("test").(int)
		s.Set("test", c+1)
		fmt.Fprintf(w, "%d", c)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	assert.Equal(t, "0", w.Body.String())

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, "1", w.Body.String())
	assert.Equal(t, 2, setCalled)
}

func TestSecureFlag(t *testing.T) {
	cases := []struct {
		tls      bool
		flag     session.Secure
		expected bool
	}{
		{false, session.NoSecure, false},
		{false, session.ForceSecure, true},
		{false, session.PreferSecure, false},
		{true, session.NoSecure, false},
		{true, session.ForceSecure, true},
		{true, session.PreferSecure, true},
	}

	for _, c := range cases {
		h := session.Middleware(session.Config{
			Store:  &mockStore{},
			Secure: c.flag,
			Proxy:  true,
		})(mockHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		if c.tls {
			r.Header.Set("X-Forwarded-Proto", "https")
		}
		h.ServeHTTP(w, r)

		cs := w.Result().Cookies()
		assert.Len(t, cs, 1)
		assert.Equal(t, c.expected, cs[0].Secure)
	}
}

func TestSecureFlagWithoutProxy(t *testing.T) {
	cases := []struct {
		flag     session.Secure
		expected bool
	}{
		{session.NoSecure, false},
		{session.ForceSecure, true},
		{session.PreferSecure, false},
	}

	for _, c := range cases {
		h := session.Middleware(session.Config{
			Store:  &mockStore{},
			Secure: c.flag,
			Proxy:  false,
		})(mockHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Forwarded-Proto", "https")
		h.ServeHTTP(w, r)

		cs := w.Result().Cookies()
		assert.Len(t, cs, 1)
		assert.Equal(t, c.expected, cs[0].Secure)
	}
}

func TestHttpOnlyFlag(t *testing.T) {
	cases := []struct {
		flag bool
	}{
		{false},
		{true},
	}

	for _, c := range cases {
		h := session.Middleware(session.Config{
			Store:    &mockStore{},
			HTTPOnly: c.flag,
		})(mockHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		h.ServeHTTP(w, r)

		cs := w.Result().Cookies()
		assert.Len(t, cs, 1)
		assert.Equal(t, c.flag, cs[0].HttpOnly)
	}
}

func TestSameSiteFlag(t *testing.T) {
	cases := []struct {
		flag session.SameSite
	}{
		{session.SameSiteNone},
		{session.SameSiteLax},
		{session.SameSiteStrict},
	}

	for _, c := range cases {
		h := session.Middleware(session.Config{
			Store:    &mockStore{},
			SameSite: c.flag,
		})(mockHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		h.ServeHTTP(w, r)

		cs := w.Result().Cookies()
		assert.Len(t, cs, 1)
		if c.flag == session.SameSiteNone {
			assert.Len(t, cs[0].Unparsed, 0)
		} else {
			assert.Equal(t, "SameSite="+string(c.flag), cs[0].Unparsed[0])
		}
	}
}

func TestRegenerate(t *testing.T) {
	c := 0

	var (
		setKey   string
		setValue = make(map[string]session.Data)
	)

	h := session.Middleware(session.Config{
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setValue[key] = value
				if c == 0 {
					setKey = key
					return nil
				}
				assert.NotEqual(t, setKey, key, "expected key after regenerate to renew")
				return nil
			},
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				return setValue[key], nil
			},
			DelFunc: func(key string, opt session.StoreOption) error {
				setValue[key] = nil
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else if c == 1 {
			s.Set("test", 2)

			// test regenerate multiple time should do nothing
			oldID := s.ID()
			s.Regenerate()
			newID := s.ID()
			assert.NotEqual(t, oldID, newID)
			s.Regenerate()
			assert.Equal(t, newID, s.ID())

			s.Set("test", 3)
			c = 2
		}
		fmt.Fprint(w, s.Get("test"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	assert.Equal(t, "1", w.Body.String())

	sess1 := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess1)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, "3", w.Body.String())

	sess2 := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess1)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, "2", w.Body.String())

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess2)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, "3", w.Body.String())
}

func TestRegenerateDeleteOldSession(t *testing.T) {
	c := 0
	setValue := make(map[string]session.Data)

	h := session.Middleware(session.Config{
		DeleteOldSession: true,
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setValue[key] = value
				return nil
			},
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				return setValue[key], nil
			},
			DelFunc: func(key string, opt session.StoreOption) error {
				setValue[key] = nil
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else if c == 1 {
			s.Set("test", 2)
			s.Regenerate()
			s.Set("test", 3)
			c = 2
		}
		fmt.Fprint(w, s.Get("test"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	assert.Equal(t, "1", w.Body.String())

	sess1 := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess1)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, "3", w.Body.String())

	sess2 := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess1)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, "<nil>", w.Body.String())

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess2)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, "3", w.Body.String())
}

func TestResave(t *testing.T) {
	setCalled := 0
	setValue := make(map[string]session.Data)

	h := session.Middleware(session.Config{
		Resave: true,
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setCalled++
				setValue[key] = value
				return nil
			},
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				return setValue[key], nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.Get(r.Context(), sessName)
		if setCalled == 0 {
			sess.Set("a", 1)
		}
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	sess := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	assert.Equal(t, 2, setCalled)
}

func TestRolling(t *testing.T) {
	c := 0

	h := session.Middleware(session.Config{
		MaxAge:  time.Second,
		Rolling: true,
		Store:   memory.New(memory.Config{}),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else {
			assert.Equal(t, 1, s.Get("test"))
		}
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	time.Sleep(time.Millisecond * 600)
	oldCookie := w.Result().Cookies()[0].Value

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Result().Cookies()[0].String())
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Len(t, w.Result().Cookies(), 1)
	assert.Equal(t, oldCookie, w.Result().Cookies()[0].Value)
}

func TestRollingDisable(t *testing.T) {
	c := 0

	h := session.Middleware(session.Config{
		MaxAge:  time.Second,
		Rolling: false,
		Store:   memory.New(memory.Config{}),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else {
			assert.Equal(t, 1, s.Get("test"))
			s.Set("test2", 1)
		}
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	time.Sleep(time.Millisecond * 600)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Result().Cookies()[0].String())
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Len(t, w.Result().Cookies(), 0)
}

func TestDestroy(t *testing.T) {
	c := 0

	var (
		delCalled bool
		setKey    string
		setValue  session.Data
	)

	h := session.Middleware(session.Config{
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setKey = key
				setValue = value
				return nil
			},
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				return setValue, nil
			},
			DelFunc: func(key string, opt session.StoreOption) error {
				delCalled = true
				assert.Equal(t, setKey, key, "expected destroy old key")
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else {
			s.Destroy()
		}
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.True(t, delCalled)
}

func TestDisableHashID(t *testing.T) {
	var setKey string

	h := session.Middleware(session.Config{
		DisableHashID: true,
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setKey = key
				return nil
			},
		},
	})(mockHandler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	cs := w.Result().Cookies()
	assert.Len(t, cs, 1)
	assert.Equal(t, setKey, cs[0].Value, "expected session id was not hashed")
}

func TestSessionMultipleGet(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), "sess")
		s.Set("test", 1)

		s = session.Get(r.Context(), "sess")
		assert.Equal(t, 1, s.Get("test"), "expected get session 2 times must not mutated value")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
}

func TestEmptyContext(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r, "expected get session from empty context must not panic")
	}()
	s := session.Get(context.Background(), "sess")
	assert.Nil(t, s, "expected get session from empty context returns nil")
}

func TestFlash(t *testing.T) {
	i := 0
	h := middleware.Chain(
		session.Middleware(session.Config{Store: memory.New(memory.Config{}), MaxAge: time.Minute}),
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), "sess")
		if i == 0 {
			s.Flash().Set("a", "1")
			s.Flash().Set("b", "2")
			i = 1
			w.Write(nil)
			return
		}
		assert.Equal(t, "1", s.Flash().Get("a"), "expected flash save in session")
		assert.Equal(t, "2", s.Flash().Get("b"), "expected flash save in session")
		w.Write(nil)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range resp.Result().Cookies() {
		req.AddCookie(c)
	}
	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	assert.Equal(t, 1, i)
}

func TestHijack(t *testing.T) {
	session.HijackedTime = 5 * time.Millisecond

	c := 0

	setValue := make(map[string]session.Data)

	h := session.Middleware(session.Config{
		Store: &mockStore{
			SetFunc: func(key string, value session.Data, opt session.StoreOption) error {
				setValue[key] = value
				return nil
			},
			GetFunc: func(key string, opt session.StoreOption) (session.Data, error) {
				return setValue[key], nil
			},
			DelFunc: func(key string, opt session.StoreOption) error {
				setValue[key] = nil
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else if c == 1 {
			s.Regenerate()
			s.Set("test", 2)
			c = 2
		} else if c == 2 {
			assert.True(t, s.Hijacked())
			c = 3
		} else if c == 3 {
			assert.False(t, s.Hijacked())
		}
		fmt.Fprint(w, "ok")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	sess1 := w.Header().Get("Set-Cookie")

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess1)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	sess2 := w.Header().Get("Set-Cookie")

	time.Sleep(session.HijackedTime)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess1)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", sess2)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
}

func TestSignature(t *testing.T) {
	c := 0

	store := memory.New(memory.Config{})
	hh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		if c == 0 {
			s.Set("a", 1)
			c++
		} else if c == 999 {
			assert.True(t, s.IsNew())
			assert.Nil(t, s.Get("a"))
		} else {
			assert.Equal(t, s.GetInt("a"), 1)
		}
		w.Write([]byte("ok"))
	})

	h := session.Middleware(session.Config{
		Keys: [][]byte{
			[]byte("key1"),
		},
		Store: store,
	})(hh)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	cs := w.Result().Cookies()
	assert.Len(t, cs, 1)
	assert.Contains(t, cs[0].Value, ".")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(cs[0])
	h.ServeHTTP(w, r)

	h = session.Middleware(session.Config{
		Keys: [][]byte{
			[]byte("key2"),
			[]byte("key1"),
		},
		Store:   store,
		Rolling: true,
	})(hh)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(cs[0])
	h.ServeHTTP(w, r)
	cs1 := w.Result().Cookies()

	h = session.Middleware(session.Config{
		Keys: [][]byte{
			[]byte("key2"),
		},
		Store: store,
	})(hh)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(cs1[0])
	h.ServeHTTP(w, r)

	// invalid signature
	c = 999
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(cs[0])
	h.ServeHTTP(w, r)
}

func TestEmptyBody(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: memory.New(memory.Config{}),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), sessName)
		s.Set("a", 1)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	cs := w.Result().Cookies()
	assert.Len(t, cs, 1)
}

func BenchmarkDefaultConfig(b *testing.B) {
	h := session.Middleware(session.Config{
		Store: &mockStore{},
	})(mockHandler)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
	}
}
