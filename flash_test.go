package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlash(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		f := new(Flash)

		assert.Empty(t, f.v, "data should be empty")
		assert.Zero(t, f.Count(), "count should be zero")
		assert.False(t, f.Changed(), "should not in changed state")

		b, err := f.encode()
		if assert.NoError(t, err, "should be able to encode") {
			assert.NotNil(t, b, "encoded value should not be nil")
			assert.Empty(t, b, "encoded value should be empty")
			assert.False(t, f.Changed(), "should still not in changed state")
		}
	})

	t.Run("Add", func(t *testing.T) {
		f := new(Flash)

		f.Add("a", 1)
		assert.True(t, f.Changed(), "should be in changed state")
		assert.EqualValues(t, 1, f.Count(), "count should changed")
	})

	t.Run("Clear Empty", func(t *testing.T) {
		f := new(Flash)

		f.Clear()

		assert.Empty(t, f.v, "data should be empty")
		assert.False(t, f.Changed(), "should not be in changed state")
	})

	t.Run("Clear Not Empty", func(t *testing.T) {
		f := new(Flash)
		f.Add("a", 1)

		f.Clear()

		assert.Empty(t, f.v, "data should be empty")
		assert.Zero(t, f.Count(), "count should be zero")
		assert.True(t, f.Changed(), "should in changed state")
	})

	t.Run("Clear Not Empty More than 1 time", func(t *testing.T) {
		f := new(Flash)
		f.Add("a", 1)

		f.Clear()
		f.Clear()

		assert.Empty(t, f.v, "data should be empty")
		assert.Zero(t, f.Count(), "count should be zero")
		assert.True(t, f.Changed(), "should in changed state")
	})

	t.Run("Del", func(t *testing.T) {
		f := new(Flash)
		f.Add("a", 1)

		f.Del("a")

		assert.True(t, f.Changed(), "should in changed state after delete key")
	})

	t.Run("Count", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			f := new(Flash)

			assert.EqualValues(t, f.Count(), 0)
		})

		t.Run("1 value", func(t *testing.T) {
			f := new(Flash)
			f.Set("a", true)

			assert.EqualValues(t, f.Count(), 1)
		})

		t.Run("2 values", func(t *testing.T) {
			f := new(Flash)
			f.Set("a", 1)
			f.Set("b", 2)

			assert.EqualValues(t, f.Count(), 2)
		})
	})

	t.Run("Clone", func(t *testing.T) {
		f := new(Flash)
		f.Add("a", "1")
		f.Add("a", "2")
		f.Add("b", "3")

		p := f.Clone()

		assert.NotEqual(t, f, p, "should not point to original")
		assert.Equal(t, f.Count(), p.Count(), "should have same count")
		assert.Equal(t, f.v, p.v, "should have same data value")

		f.Clear()
		assert.NotEqual(t, f.v, p.v)
	})

	t.Run("Values", func(t *testing.T) {
		f := new(Flash)

		v := f.Values("a")
		assert.NotNil(t, v, "not exists key should not nill")
		assert.Empty(t, v, "not exists key should empty")

		f.Add("a", 1)
		f.Add("a", 2)
		f.Add("a", 3)
		v = f.Values("a")
		assert.Equal(t, []interface{}{1, 2, 3}, v)
		assert.False(t, f.Has("a"), "key should deleted after call values")
	})
}

func TestFlashEncodeDecode(t *testing.T) {
	t.Parallel()

	t.Run("Valid data", func(t *testing.T) {
		f := new(Flash)
		f.Add("a", 1)

		b, err := f.encode()
		if assert.NoError(t, err) {
			assert.NotNil(t, b)
			assert.NotEmpty(t, b)
			assert.True(t, f.Changed())

			p := new(Flash)
			if assert.NoError(t, p.decode(b)) {
				assert.Equal(t, f.Count(), p.Count())
				assert.Equal(t, f.v, p.v)
			}
		}
	})

	t.Run("Invalid data", func(t *testing.T) {
		f := new(Flash)
		f.Set("key", &struct{}{})

		b, err := f.encode()
		assert.Error(t, err)
		assert.Empty(t, b)
	})

	t.Run("Decode empty bytes", func(t *testing.T) {
		f := new(Flash)

		if assert.NoError(t, f.decode([]byte{})) {
			assert.Empty(t, f.v)
			assert.False(t, f.Changed())
		}
	})

	t.Run("Decode invalid bytes", func(t *testing.T) {
		f := new(Flash)

		if assert.Error(t, f.decode([]byte("invalid data"))) {
			assert.Empty(t, f.v)
			assert.False(t, f.Changed())
		}
	})
}

func TestFlashGet(t *testing.T) {
	t.Parallel()

	t.Run("Get empty", func(t *testing.T) {
		f := new(Flash)

		assert.Nil(t, f.Get("a"))
		assert.False(t, f.Changed())
	})

	t.Run("Get valid value", func(t *testing.T) {
		f := new(Flash)
		f.Set("a", 1)

		assert.True(t, f.Has("a"))
		assert.Equal(t, 1, f.Get("a"))
		assert.False(t, f.Has("a"))
	})

	t.Run("GetString from string value", func(t *testing.T) {
		f := new(Flash)
		f.Set("a", "hello")

		assert.Equal(t, "hello", f.GetString("a"))
	})

	t.Run("GetInt from string value", func(t *testing.T) {
		f := new(Flash)
		f.Set("a", "hello")

		assert.Zero(t, f.GetInt("a"))
	})

	t.Run("GetInt64 from string value", func(t *testing.T) {
		f := new(Flash)
		f.Set("a", "hello")

		assert.Zero(t, f.GetInt64("a"))
	})

	t.Run("GetFloat32 from string value", func(t *testing.T) {
		f := new(Flash)
		f.Set("a", "hello")

		assert.Zero(t, f.GetFloat32("a"))
	})

	t.Run("GetFloat64 from string value", func(t *testing.T) {
		f := new(Flash)
		f.Set("a", "hello")

		assert.Zero(t, f.GetFloat64("a"))
	})

	t.Run("GetBool from string value", func(t *testing.T) {
		f := new(Flash)
		f.Set("a", "hello")

		assert.Zero(t, f.GetBool("a"))
	})
}
