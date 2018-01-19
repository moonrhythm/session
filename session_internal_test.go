package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
