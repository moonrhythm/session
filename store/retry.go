package store

import (
	"context"
	"errors"
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
func (s *Retry) Get(ctx context.Context, key string) (r session.Data, err error) {
	for range s.maxAttempts() {
		r, err = s.Store.Get(ctx, key)
		if err == nil || errors.Is(err, session.ErrNotFound) {
			break
		}
		time.Sleep(s.backOffDuration())
	}
	return
}

// Set sets session data to wrapped store with retry
func (s *Retry) Set(ctx context.Context, key string, value session.Data, opt session.StoreOption) (err error) {
	for range s.maxAttempts() {
		err = s.Store.Set(ctx, key, value, opt)
		if err == nil {
			break
		}
		time.Sleep(s.backOffDuration())
	}
	return
}

// Del deletes session data from wrapped store with retry
func (s *Retry) Del(ctx context.Context, key string) (err error) {
	for range s.maxAttempts() {
		err = s.Store.Del(ctx, key)
		if err == nil {
			break
		}
		time.Sleep(s.backOffDuration())
	}
	return
}
