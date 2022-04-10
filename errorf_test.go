package whoops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err1 := Errorf("hello")

	err := err1.Format()
	t.Run("Errorf generates an error", func(t *testing.T) {
		_, ok := any(err).(error)
		assert.True(t, ok)
	})

	t.Run("testing Is", func(t *testing.T) {
		var (
			errf1 = Errorf("err: %s")
			errf2 = Errorf("err: %s, %d")
		)

		var (
			err1  = errf1.Format("bla")
			err21 = errf2.Format("bar", 3)
			err22 = errf2.Format("foo", 8, 9)
		)

		assert.False(t, err1.Is(err21))
		assert.False(t, err1.Is(err22))
		assert.False(t, err21.Is(err22))
		assert.False(t, err22.Is(err21))

		assert.True(t, err1.Is(err1))
		assert.True(t, err21.Is(err21))
		assert.True(t, err22.Is(err22))
	})

	t.Run("formatting", func(t *testing.T) {
		var (
			errf = Errorf("err: %s, %d")
			err  = errf.Format("bob", 1337)
		)

		assert.Equal(t, "err: bob, 1337", err.Error())
	})
}
