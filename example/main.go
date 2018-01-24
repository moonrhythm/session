package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/acoshift/session"
	"github.com/acoshift/session/store/memory"
)

func main() {
	h := session.New(session.Config{
		Store:    memory.New(memory.Config{}),
		HTTPOnly: true,
		Secret:   []byte("supersalt"),
		Keys:     [][]byte{[]byte("supersecret")},
		Path:     "/",
		Rolling:  true,
		MaxAge:   time.Hour,
		SameSite: session.SameSiteLax,
		Secure:   session.PreferSecure,
		Proxy:    true,
	}).Middleware()(http.HandlerFunc(handler))
	http.ListenAndServe(":8080", h)
}

func handler(w http.ResponseWriter, r *http.Request) {
	sess := session.Get(r.Context(), "sess")
	cnt := sess.GetInt("cnt")
	cnt++
	sess.Set("cnt", cnt)
	fmt.Fprintf(w, "%d views", cnt)
}
