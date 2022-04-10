package whoops

import "github.com/stretchr/testify/assert"
import "testing"
import "errors"

type testStruct struct {
	A string
	B int
}

func TestFields(t *testing.T) {
	var (
		f1 Field[string]     = "field1"
		f2 Field[testStruct] = "field2"
		f3 Field[int]        = "field which is not used"
		f4 Field[bool]       = "something"

		ts = testStruct{
			A: "A",
			B: 2,
		}
	)

	anyErr := errors.New("any error")

	we := Wrap(anyErr, f1.Val("random value"), f2.Val(ts))
	_, isErr := any(we).(error)

	assert.True(t, isErr)

	extractedField1, found := f1.GetFrom(we)
	assert.True(t, found)
	assert.Equal(t, "random value", extractedField1)

	extractedField2, found := f2.GetFrom(we)
	assert.True(t, found)
	assert.Equal(t, ts, extractedField2)

	_, found = f3.GetFrom(we)
	assert.False(t, found)

	newWe := Wrap(we, f4.Val(true))

	// f1, f2 and f4 must exist on newWe

	var ok bool
	{
		_, ok = f1.GetFrom(newWe)
		assert.True(t, ok)
		_, ok = f2.GetFrom(newWe)
		assert.True(t, ok)
		_, ok = f4.GetFrom(newWe)
		assert.True(t, ok)
	}

	// f4 does not exist on old we

	v, ok := f4.GetFrom(we)
	assert.False(t, ok)
	assert.False(t, v)

}
