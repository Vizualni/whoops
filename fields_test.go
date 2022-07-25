package whoops

import "github.com/stretchr/testify/assert"
import "testing"
import "errors"

type testStruct struct {
	A string
	B int
}

func TestErrorEnrichingWithFields(t *testing.T) {
	var (
		f1 Field[string]     = "field1"
		f2 Field[testStruct] = "field2"
		f3 Field[int]        = "field which is not used"
		f4 Field[bool]       = "something"

		ts = testStruct{
			A: "A",
			B: 2,
		}

		enrichedErr enrichedErrorWithFields
	)

	anyErr := errors.New("any error")

	t.Run("Enrich enriches an error", func(t *testing.T) {
		enrichedErr = Enrich(anyErr, f1.Val("random value"), f2.Val(ts))
		_, isErr := any(enrichedErr).(error)
		assert.True(t, isErr)
		assert.Equal(t, "original error: any error\nenriched fields:\nfield1[string] = \"random value\"\nfield2[whoops.testStruct] = whoops.testStruct{A:\"A\", B:2}\n", enrichedErr.Error())
	})

	t.Run("ensure that fields can get extracted", func(t *testing.T) {
		extractedField1, found := f1.GetFrom(enrichedErr)
		assert.True(t, found)
		assert.Equal(t, "random value", extractedField1)

		extractedField2, found := f2.GetFrom(enrichedErr)
		assert.True(t, found)
		assert.Equal(t, ts, extractedField2)

		_, found = f3.GetFrom(enrichedErr)
		assert.False(t, found)
	})

	newEnriched := Enrich(enrichedErr, f4.Val(true))

	t.Run("enriching an enriched error", func(t *testing.T) {
		// f1, f2 and f4 must exist on newWe
		var ok bool
		_, ok = f1.GetFrom(newEnriched)
		assert.True(t, ok)
		_, ok = f2.GetFrom(newEnriched)
		assert.True(t, ok)
		_, ok = f4.GetFrom(newEnriched)
		assert.True(t, ok)

		// f4 does not exist on old we
		v, ok := f4.GetFrom(enrichedErr)
		assert.False(t, ok)
		assert.False(t, v)
	})

}
