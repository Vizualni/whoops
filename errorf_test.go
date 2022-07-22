package whoops

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorf(t *testing.T) {
	t.Run("Errorf generates an error", func(t *testing.T) {
		err := Errorf("hello").Format()
		_, ok := any(err).(error)
		require.True(t, ok)
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

		require.NotErrorIs(t, err1, err21)
		require.NotErrorIs(t, err21, err1)

		require.NotErrorIs(t, err1, err22)
		require.NotErrorIs(t, err22, err1)

		require.NotErrorIs(t, err21, err22)
		require.NotErrorIs(t, err22, err21)

		require.ErrorIs(t, err1, err1)
		require.ErrorIs(t, err21, errf2)
		require.ErrorIs(t, errf2, err21)

		require.ErrorIs(t, err22, err22)
	})

	t.Run("formatting", func(t *testing.T) {
		var (
			errf = Errorf("err: %s, %d")
			err  = errf.Format("bob", 1337)
		)

		require.Equal(t, "err: bob, 1337", err.Error())
	})
}
