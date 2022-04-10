package whoops

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err1 := Errorf("hello")

	err := err1.Format()
	t.Run("formatting errorf returns an error", func(t *testing.T) {
		_, ok := any(err).(error)
		assert.True(t, ok)
	})

	t.Run("testing Is", func(t *testing.T) {
		t.Run("err is err", func(t *testing.T) {
			assert.True(t, err.Is(err))
		})

		t.Run("err is not different error", func(t *testing.T) {
			assert.False(t, err.Is(errors.New("different error")))
		})

		t.Run("err is Errorf", func(t *testing.T) {
			assert.True(t, err.Is(err1))
		})

		t.Run("err is not another Errorf", func(t *testing.T) {
			assert.False(t, err.Is(Errorf("different error")))
		})
	})
}

func TestErrorfIs(t *testing.T) {
	var (
		errf1 = Errorf("err: %s")
		errf2 = Errorf("err: %s, %d")
	)

	var (
		err1  = errf1.Format("bla")
		err21 = errf2.Format("bar", 3)
		err22 = errf2.Format("foo", 8)
	)

	assert.False(t, err1.Is(err21))
	assert.False(t, err1.Is(err22))
	assert.False(t, err21.Is(err22))
	assert.False(t, err22.Is(err21))

	assert.True(t, err1.Is(err1))
	assert.True(t, err21.Is(err21))
	assert.True(t, err22.Is(err22))
}

func TestErrorfFormatting(t *testing.T) {
	var (
		errf = Errorf("err: %s, %d")
		err  = errf.Format("bob", 1337)
	)

	assert.Equal(t, "err: bob, 1337", err.Error())
}
