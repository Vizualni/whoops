package whoops

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	var g Group
	g.Add(nil)
	assert.False(t, g.Err())
	assert.Equal(t, "no errors in the group", g.Error())

	g.Add(errors.New("something new1"))
	g.Add(errors.New("something new2"))
	assert.True(t, g.Err())

	// not going to assert the whole string, just that it starts with "there are 2 errors in the group"
	assert.True(t, strings.HasPrefix(g.Error(), "there are 2 errors in the group"), g.Error())
}

func TestGroupInitWithMultipleErrors(t *testing.T) {
	err := Group{
		errors.New("a"),
		errors.New("b"),
	}
	assert.True(t, err.Err())
	err = Group{}
	err.Add(
		errors.New("c"),
		errors.New("d"),
	)
	assert.True(t, err.Err())
}

func TestGroupUnwrapping(t *testing.T) {
	var (
		err1 = errors.New("err1")
		err2 = errors.New("err2")
		err3 = errors.New("err3")
		err4 = errors.New("err3")
	)
	err := Group{err1, err2}
	assert.True(t, Is(err, err1))
	assert.True(t, Is(err, err2))
	assert.False(t, Is(err, err3))

	unwrappedErr := Unwrap(err)
	assert.True(t, Is(unwrappedErr, err1))
	unwrappedErr = Unwrap(unwrappedErr)
	assert.True(t, Is(unwrappedErr, err2))
	unwrappedErr = Unwrap(unwrappedErr)
	assert.Nil(t, unwrappedErr)

	t.Run("two groups unwrapping", func(t *testing.T) {
		var (
			g1 = Group{err1, err2}
			g2 = Group{err3, err4}
			g3 = Group{g1, g2}
		)
		for i, err := range []error{err1, err2, err3, err4} {
			t.Run(fmt.Sprintf("err index %d", i), func(t *testing.T) {
				assert.True(t, Is(g3, err))
			})
		}

	})

}
