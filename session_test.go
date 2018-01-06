package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeEmpty(t *testing.T) {
	s := Session{}
	b, err := s.MarshalBinary()
	assert.NoError(t, err)
	assert.NotNil(t, b, "expected encode always return not nil")
	assert.Len(t, b, 0)
}

func TestEncodeUnregisterType(t *testing.T) {
	type a struct{}
	s := Session{}
	s.Set("a", a{})
	b, err := s.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestSessionOperation(t *testing.T) {
	s := Session{}
	assert.Nil(t, s.Get("a"), "expected get data from empty session return nil")
	assert.Nil(t, s.Pop("a"), "expected pop data from empty session return nil")

	s.Del("a")
	assert.Nil(t, s.data)

	s.Set("a", 1)
	assert.Equal(t, 1, s.Get("a"))

	s.Del("a")
	assert.Nil(t, s.Get("a"), "expected get data after delete to be nil")

	s.Set("b", 1)
	assert.Equal(t, 1, s.Pop("b"))
	assert.Nil(t, s.Get("b"))
}

func TestRenew(t *testing.T) {
	s := Session{}
	s.Set("a", "1")
	s.Renew()
	assert.Nil(t, s.Get("a"), "expected renew must delete all data")
}
