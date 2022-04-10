package whoops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapper(t *testing.T) {
	var (
		err1 = String("1")
		err2 = String("2")
		w    = Wrap(err1, err2)
	)

	assert.True(t, Is(w, err1))
	assert.True(t, Is(w, err2))
	assert.True(t, Is(w, w))
}
