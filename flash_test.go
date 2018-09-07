package session

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFlashNew(t *testing.T) {
	Convey("Given new flash", t, func() {
		f := new(Flash)

		Convey("Flash should not be nil", func() {
			So(f, ShouldNotBeNil)
		})

		Convey("Flash data should be empty", func() {
			So(f.v, ShouldBeEmpty)
		})

		Convey("Flash count should be empty", func() {
			So(f.Count(), ShouldEqual, 0)
		})

		Convey("Flash should not in changed state", func() {
			So(f.Changed(), ShouldBeFalse)
		})

		Convey("Flash should be able to encode", func() {
			b, err := f.encode()

			Convey("Without error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Encoded result should not be nil", func() {
				So(b, ShouldNotBeNil)
			})

			Convey("Encoded bytes should be empty", func() {
				So(b, ShouldBeEmpty)
			})

			Convey("Flash still unchanged", func() {
				So(f.Changed(), ShouldBeFalse)
			})
		})

		Convey("When add a data to flash", func() {
			f.Add("a", 1)

			Convey("Flash should be in changed state", func() {
				So(f.Changed(), ShouldBeTrue)
			})

			Convey("Count should changed", func() {
				So(f.Count(), ShouldEqual, 1)
			})
		})
	})
}

func TestFlashClear(t *testing.T) {
	Convey("Given empty flash", t, func() {
		f := new(Flash)

		Convey("When clear", func() {
			f.Clear()

			Convey("Should not be in changed state", func() {
				So(f.Changed(), ShouldBeFalse)
			})
		})
	})

	Convey("Given not empty flash", t, func() {
		f := new(Flash)
		f.Add("a", 1)

		Convey("When clear", func() {
			f.Clear()

			Convey("Flash data should be empty", func() {
				So(f.v, ShouldBeEmpty)
			})

			Convey("Count should be zero", func() {
				So(f.Count(), ShouldBeZeroValue)
			})

			Convey("Should be in changed state", func() {
				So(f.Changed(), ShouldBeTrue)
			})

			Convey("When clear again", func() {
				Convey("Should still be in changed state", func() {
					So(f.Changed(), ShouldBeTrue)
				})
			})
		})
	})
}

func TestFlashEncodeDecode(t *testing.T) {
	Convey("Given not empty flash", t, func() {
		f := new(Flash)
		f.Add("a", 1)

		Convey("Should be able to encode", func() {
			b, err := f.encode()

			Convey("Without error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Encoded result should not be nil", func() {
				So(b, ShouldNotBeNil)
			})

			Convey("Encoded bytes should not be empty", func() {
				So(b, ShouldNotBeEmpty)
			})

			Convey("Flash should still be changed state", func() {
				So(f.Changed(), ShouldBeTrue)
			})

			Convey("Encoded data should be able to decode", func() {
				f := new(Flash)
				err := f.decode(b)

				Convey("Without error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Decoded flash should not be nil", func() {
					So(f, ShouldNotBeNil)
				})

				Convey("With same count", func() {
					So(f.Count(), ShouldEqual, 1)
				})

				Convey("With same data", func() {
					So(f.v["a"], ShouldNotBeEmpty)
					So(f.v["a"][0], ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given flash with invalid data", t, func() {
		f := new(Flash)
		f.Set("key", &struct{}{})

		Convey("When encode", func() {
			b, err := f.encode()

			Convey("Should error", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Encoded bytes should be empty", func() {
				So(b, ShouldBeEmpty)
			})
		})
	})

	Convey("Given empty bytes", t, func() {
		b := []byte{}

		Convey("Bytes should be able to decode", func() {
			f := new(Flash)
			err := f.decode(b)

			Convey("Without error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Flash should not be nil", func() {
				So(f, ShouldNotBeNil)
			})

			Convey("Flash data should be empty", func() {
				So(f.v, ShouldBeEmpty)
			})

			Convey("Flash count should be zero", func() {
				So(f.Count(), ShouldBeZeroValue)
			})

			Convey("Flash should not in changed state", func() {
				So(f.Changed(), ShouldBeFalse)
			})
		})
	})

	Convey("Given invalid encoded flash bytes", t, func() {
		b := []byte("invalid data")

		Convey("When decode", func() {
			f := new(Flash)
			err := f.decode(b)

			Convey("Should error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestFlashDel(t *testing.T) {
	Convey("Given not empty flash", t, func() {
		f := new(Flash)
		f.Add("a", 1)

		Convey("When delete all keys", func() {
			f.Del("a")

			Convey("Should still be in changed state", func() {
				So(f.Changed(), ShouldBeTrue)
			})
		})
	})
}

func TestFlashValues(t *testing.T) {
	Convey("Given empty flash", t, func() {
		f := new(Flash)

		Convey("When retrieve values from not exists key", func() {
			v := f.Values("a")

			Convey("Should return non-nil values", func() {
				So(v, ShouldNotBeNil)
			})

			Convey("Should return empty data", func() {
				So(v, ShouldBeEmpty)
			})
		})

		Convey("When add 3 values same key", func() {
			f.Add("a", 1)
			f.Add("a", 2)
			f.Add("a", 3)

			Convey("When retrieve values for that key", func() {
				v := f.Values("a")

				Convey("Should resemble added values", func() {
					So(v, ShouldResemble, []interface{}{1, 2, 3})
				})

				Convey("Should not disapper in flash", func() {
					So(f.Has("a"), ShouldBeFalse)
				})
			})
		})
	})
}

func TestFlashGet(t *testing.T) {
	Convey("Given empty flash", t, func() {
		f := new(Flash)

		Convey("When get not exists key", func() {
			v := f.Get("a")

			Convey("Should return nil", func() {
				So(v, ShouldBeNil)
			})

			Convey("Should still be unchanged state", func() {
				So(f.Changed(), ShouldBeFalse)
			})
		})

		Convey("When set a value", func() {
			f.Set("a", 1)

			Convey("Should has that key", func() {
				So(f.Has("a"), ShouldBeTrue)
			})

			Convey("When get that key", func() {
				v := f.Get("a")

				Convey("Should return same value", func() {
					So(v, ShouldEqual, 1)
				})

				Convey("Key should disappear", func() {
					So(f.Has("a"), ShouldBeFalse)
				})
			})
		})

		Convey("When set a string value", func() {
			f.Set("a", "hello")

			Convey("When get string from that key", func() {
				v := f.GetString("a")

				Convey("Should return same value", func() {
					So(v, ShouldEqual, "hello")
				})
			})

			Convey("When get int from that key", func() {
				v := f.GetInt("a")

				Convey("Should return zero value", func() {
					So(v, ShouldBeZeroValue)
				})
			})

			Convey("When get int64 from that key", func() {
				v := f.GetInt64("a")

				Convey("Should return zero value", func() {
					So(v, ShouldBeZeroValue)
				})
			})

			Convey("When get float32 from that key", func() {
				v := f.GetFloat32("a")

				Convey("Should return zero value", func() {
					So(v, ShouldBeZeroValue)
				})
			})

			Convey("When get float64 from that key", func() {
				v := f.GetFloat64("a")

				Convey("Should return zero value", func() {
					So(v, ShouldBeZeroValue)
				})
			})

			Convey("When get bool from that key", func() {
				v := f.GetBool("a")

				Convey("Should return zero value", func() {
					So(v, ShouldBeZeroValue)
				})
			})
		})
	})
}

func TestFlashClone(t *testing.T) {
	Convey("Given not empty flash", t, func() {
		f := new(Flash)
		f.Add("a", "1")
		f.Add("a", "2")
		f.Add("b", "3")

		Convey("When clone", func() {
			p := f.Clone()

			Convey("Should not point to original", func() {
				So(p, ShouldNotPointTo, f)
			})

			Convey("Data should have same count", func() {
				So(p.Count(), ShouldEqual, f.Count())
			})

			Convey("Data should have same values", func() {
				So(p.v, ShouldResemble, f.v)
			})

			Convey("When clear original flash", func() {
				f.Clear()

				Convey("Clone data should not resemble original data", func() {
					So(p.v, ShouldNotResemble, f.v)
				})
			})
		})

	})
}

func TestFlashCount(t *testing.T) {
	Convey("Given empty flash", t, func() {
		f := new(Flash)

		Convey("Should have 0 count", func() {
			So(f.Count(), ShouldEqual, 0)
		})

		Convey("When set a value", func() {
			f.Set("a", true)

			Convey("Should have 1 count", func() {
				So(f.Count(), ShouldEqual, 1)
			})

			Convey("When get that value", func() {
				f.Get("a")

				Convey("Should back to 0 count", func() {
					So(f.Count(), ShouldEqual, 0)
				})
			})
		})

		Convey("When set 2 values", func() {
			f.Set("a", 1)
			f.Set("b", 2)

			Convey("Should have 2 count", func() {
				So(f.Count(), ShouldEqual, 2)
			})
		})
	})
}
