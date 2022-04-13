package whoops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatch(t *testing.T) {
	const err1 = String("ohai")
	t.Run("catch returns the returned error", func(t *testing.T) {
		err := Catch(func() error {
			return err1
		})
		assert.ErrorIs(t, err, err1)
	})

	t.Run("when calling must from catch it returns that error", func(t *testing.T) {
		err := Catch(func() error {
			f1 := func() error {
				return err1
			}
			Must1(f1())
			return nil
		})
		assert.ErrorIs(t, err, err1)
	})

	t.Run("when calling assert it returns the error", func(t *testing.T) {
		err := Catch(func() error {
			Assert(err1)
			return nil
		})
		assert.ErrorIs(t, err, err1)
	})

	t.Run("when something else panics, it raises it", func(t *testing.T) {
		assert.PanicsWithValue(t, err1, func() {
			Catch(func() error {
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
