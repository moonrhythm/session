package retry

import (
	"github.com/acoshift/session"
)

// New creates new retry store
func New(store session.Store, maxAttempts int) session.Store {
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	return &retryStore{
		store:       store,
		maxAttempts: maxAttempts,
	}
}

type retryStore struct {
	store       session.Store
	maxAttempts int
}

func (s *retryStore) Get(key string, opt session.StoreOption) (r session.Data, err error) {
	for i := 0; i < s.maxAttempts; i++ {
		r, err = s.store.Get(key, opt)
		if err == nil || err == session.ErrNotFound {
			break
		}
	}
	return
}

func (s *retryStore) Set(key string, value session.Data, opt session.StoreOption) (err error) {
	for i := 0; i < s.maxAttempts; i++ {
		err = s.store.Set(key, value, opt)
		if err == nil {
			break
		}
	}
	return
}

func (s *retryStore) Del(key string, opt session.StoreOption) (err error) {
	for i := 0; i < s.maxAttempts; i++ {
		err = s.store.Del(key, opt)
		if err == nil {
			break
		}
	}
	return
}
