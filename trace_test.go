package whoops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	t.Run("one level deep", func(t *testing.T) {
		err := String("something")
		terr := Trace(err)
		assert.Len(t, FormatStacktrace(err), 0)
		assert.Len(t, FormatStacktrace(terr), 3)
	})
}

func TestTraceString(t *testing.T) {
	t.Run("one level deep", func(t *testing.T) {
		terr := String("something").Trace()
		assert.Len(t, FormatStacktrace(terr), 3)
	})
}

func TestTraceErrorf(t *testing.T) {
	t.Run("one level deep", func(t *testing.T) {
		terr := Errorf("something").Format().Trace()
		assert.Len(t, FormatStacktrace(terr), 3)
	})
}

func TestTraceGroup(t *testing.T) {
	t.Run("one level deep", func(t *testing.T) {
		terr := Group{String("bla")}.Trace()
		assert.Len(t, FormatStacktrace(terr), 3)
	})
}

func TestWrapperTrace(t *testing.T) {
	t.Run("one level deep", func(t *testing.T) {
		terr := Wrap(String("bla"), String("foo")).Trace()
		assert.Len(t, FormatStacktrace(terr), 3)
	})
}
