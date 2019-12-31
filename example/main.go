package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/moonrhythm/session"
	"github.com/moonrhythm/session/store"
)

func main() {
	h := session.New(session.Config{
		Store:    new(store.Memory),
		HTTPOnly: true,
		Secret:   []byte("supersalt"),
		Keys:     [][]byte{[]byte("supersecret")},
		Path:     "/",
		Rolling:  true,
		MaxAge:   time.Hour,
		SameSite: http.SameSiteLaxMode,
		Secure:   session.PreferSecure,
		Proxy:    true,
	}).Middleware()(http.HandlerFunc(handler))
	http.ListenAndServe(":8080", h)
}

func handler(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.Get(r.Context(), "sess")
	cnt := sess.GetInt("cnt")
	cnt++
	sess.Set("cnt", cnt)
	fmt.Fprintf(w, "%d views", cnt)
}
