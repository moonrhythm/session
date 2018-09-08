package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionOperation(t *testing.T) {
	t.Parallel()

	s := Session{}
	assert.Nil(t, s.Get("a"), "expected get data from empty session return nil")
	assert.Nil(t, s.Pop("a"), "expected pop data from empty session return nil")
	assert.False(t, s.changed, "expected pop empty key not trigger changed")

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

func TestSessionGetWithTypes(t *testing.T) {
	t.Parallel()

	s := Session{}

	s.Set("string", "text")
	s.Set("int", 10)
	s.Set("int64", int64(10))
	s.Set("true", true)
	s.Set("false", false)
	s.Set("float32", float32(1.3))
	s.Set("float64", float64(1.5))

	assert.Equal(t, s.Get("string"), s.GetString("string"))
	assert.Equal(t, s.Get("int"), s.GetInt("int"))
	assert.Equal(t, s.Get("int64"), s.GetInt64("int64"))
	assert.Equal(t, s.Get("true"), s.GetBool("true"))
	assert.Equal(t, s.Get("false"), s.GetBool("false"))
	assert.Equal(t, s.Get("float32"), s.GetFloat32("float32"))
	assert.Equal(t, s.Get("float64"), s.GetFloat64("float64"))

	assert.Equal(t, s.Get("string"), s.PopString("string"))
	assert.Equal(t, s.Get("int"), s.PopInt("int"))
	assert.Equal(t, s.Get("int64"), s.PopInt64("int64"))
	assert.Equal(t, s.Get("true"), s.PopBool("true"))
	assert.Equal(t, s.Get("false"), s.PopBool("false"))
	assert.Equal(t, s.Get("float32"), s.PopFloat32("float32"))
	assert.Equal(t, s.Get("float64"), s.PopFloat64("float64"))

	assert.Empty(t, s.PopString("string"))
	assert.Empty(t, s.PopInt("int"))
	assert.Empty(t, s.PopInt64("int64"))
	assert.Empty(t, s.PopBool("true"))
	assert.Empty(t, s.PopBool("false"))
	assert.Empty(t, s.PopFloat32("float32"))
	assert.Empty(t, s.PopFloat64("float64"))
}
