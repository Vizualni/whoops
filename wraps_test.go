package whoops

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWraps(t *testing.T) {
	err := String("oh no")
	w := WrapS(err, "abc: %d %d", 12, 34)
	require.ErrorIs(t, w, err)

	require.Equal(t, "oh no: abc: 12 34", w.Error())
}

func TestWrapsWithOtherErrors(t *testing.T) {
	e1 := String("oh no")
	e2 := Errorf("oh no").Format()
	e3 := Group{e1, e2}

	require.ErrorIs(t, e1.WrapS("abc"), e1)
	require.ErrorIs(t, e2.WrapS("abc"), e2)
	require.ErrorIs(t, e3.WrapS("abc"), e1)
}
