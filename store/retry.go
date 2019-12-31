package store

import (
	"time"

	"github.com/moonrhythm/session"
)

// Retry reties store operation when failed
type Retry struct {
	Store       session.Store
	MaxAttempts int
}

func (s *Retry) maxAttempts() int {
	if s.MaxAttempts <= 0 {
		return 3
	}
	return s.MaxAttempts
}

func (s *Retry) backOffDuration() time.Duration {
	return 100 * time.Millisecond
}

// Get gets session data from wrapped store with retry
func (s *Retry) Get(key string) (r session.Data, err error) {
	for i := 0; i < s.MaxAttempts; i++ {
		r, err = s.Store.Get(key)
		if err == nil || err == session.ErrNotFound {
			break
		}
		time.Sleep(s.backOffDuration())
	}
	return
}

// Set sets session data to wrapped store with retry
func (s *Retry) Set(key string, value session.Data, opt session.StoreOption) (err error) {
	for i := 0; i < s.MaxAttempts; i++ {
		err = s.Store.Set(key, value, opt)
		if err == nil {
			break
		}
		time.Sleep(s.backOffDuration())
	}
	return
}

// Del deletes session data from wrapped store with retry
func (s *Retry) Del(key string) (err error) {
	for i := 0; i < s.MaxAttempts; i++ {
		err = s.Store.Del(key)
		if err == nil {
			break
		}
		time.Sleep(s.backOffDuration())
	}
	return
}
