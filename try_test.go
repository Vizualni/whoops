package whoops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTry(t *testing.T) {
	const err1 = String("ohai")
	const wrapperErr = String("wrapper")

	t.Run("when calling must from try, it returns that error", func(t *testing.T) {
		err := Try(func() {
			f1 := func() error {
				return err1
			}
			Must1(f1())
		})
		assert.ErrorIs(t, err, err1)
	})

	t.Run("when calling assert it returns the error", func(t *testing.T) {
		err := Try(func() {
			Assert(err1)
		})
		assert.ErrorIs(t, err, err1)
	})

	t.Run("when something else panics, it raises it", func(t *testing.T) {
		assert.PanicsWithValue(t, err1, func() {
			Try(func() {
				panic(err1)
			})
		})
	})

	t.Run("calling defer wrap wraps the error", func(t *testing.T) {
		err := Try(func() {
			defer DeferWrap(wrapperErr)()
			Assert(err1)
		})
		assert.ErrorIs(t, err, err1)
		assert.ErrorIs(t, err, wrapperErr)
	})

	t.Run("calling defer wrap with no error skips it", func(t *testing.T) {
		err := Try(func() {
			defer DeferWrap(wrapperErr)()
		})
		assert.NoError(t, err)
	})

	t.Run("calling defer with custom panic, re-panics the same error", func(t *testing.T) {
		assert.PanicsWithValue(t, "omg", func() {
			Try(func() {
				defer DeferWrap(wrapperErr)()
				panic("omg")
			})
		})
	})
}

func TestTryVal(t *testing.T) {
	const err1 = String("ohai")

	t.Run("no error, it returns a value", func(t *testing.T) {
		val, err := TryVal(func() int {
			return 1
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
	})

	t.Run("when calling must from try, it returns that error", func(t *testing.T) {
		val, err := TryVal(func() int {
			f1 := func() error {
				return err1
			}
			Must1(f1())
			return 123 // never gets here
		})
		assert.ErrorIs(t, err, err1)
		assert.Equal(t, 0, val)
	})

	t.Run("when calling assert it returns the error", func(t *testing.T) {
		val, err := TryVal(func() int {
			Assert(err1)
			return 0
		})
		assert.ErrorIs(t, err, err1)
		assert.Equal(t, 0, val)
	})

	t.Run("when something else panics, it raises it", func(t *testing.T) {
		assert.PanicsWithValue(t, err1, func() {
			TryVal(func() bool {
				panic(err1)
			})
		})
	})
}

func TestMust1(t *testing.T) {
	f1 := func() error {
		return nil
	}
	f2 := func() error {
		return String("hello")
	}

	Must1(f1())

	assert.Panics(t, func() {
		Must1(f2())
		assert.False(t, true, "this should not have happened")
	})
}

func TestMust2(t *testing.T) {
	f1 := func() (int, error) {
		return 1, nil
	}
	f2 := func() (int, error) {
		return 0, String("hello")
	}

	val := Must(f1())
	assert.Equal(t, 1, val)

	assert.Panics(t, func() {
		Must(f2())
		assert.False(t, true, "this should not have happened")
	})
}

func TestMust3(t *testing.T) {
	f1 := func() (int, string, error) {
		return 1, "oh no", nil
	}
	f2 := func() (int, string, error) {
		return 0, "", String("hello")
	}

	val1, val2 := Must3(f1())
	assert.Equal(t, 1, val1)
	assert.Equal(t, "oh no", val2)

	assert.Panics(t, func() {
		Must3(f2())
		assert.False(t, true, "this should not have happened")
	})
}
