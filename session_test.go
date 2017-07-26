package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/acoshift/session"
)

func TestPanicConfig(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("expected panic when misconfig; but not")
		}
	}()
	session.Middleware(session.Config{})
}

func TestDefautConfig(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		s.Set("test", 1)
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	cookie := w.Header().Get("Set-Cookie")
	if len(cookie) == 0 {
		t.Fatalf("expected cookie not empty; got empty")
	}
}

func TestSessionSetInStore(t *testing.T) {
	var (
		setCalled bool
		setKey    string
		setValue  []byte
		setTTL    time.Duration
	)

	h := session.Middleware(session.Config{
		Name:   "sess",
		MaxAge: time.Second,
		Store: &mockStore{
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				setCalled = true
				setKey = key
				setValue = value
				setTTL = ttl
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		s.Set("test", 1)
		w.Write([]byte("ok"))
	}))

	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("http get error; %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status code 200; got %d", resp.StatusCode)
	}
	if !setCalled {
		t.Fatalf("expected store was called; but not")
	}
	if len(setKey) == 0 {
		t.Fatalf("expected key not empty; got empty")
	}
	if len(setValue) == 0 {
		t.Fatalf("expected value not empty; got empty")
	}
	if setTTL != time.Second {
		t.Fatalf("expected ttl to be 1s; got %v", setTTL)
	}
}
