package whoops

import (
	"fmt"
	"strings"
)

type Field[T any] string // adding T for type safety

var mm Field[int] = "bla"

type wrappedValue[T any] struct {
	typ Field[T]
	val T
}

func (f wrappedValue[T]) String() string {
	var a any = f.val
	return fmt.Sprintf("%s[%T] = %#v", f.typ, f.val, a)
}

func (f Field[T]) Val(v T) wrappedValue[T] {
	return wrappedValue[T]{
		typ: f,
		val: v,
	}
}

func (f Field[T]) GetFrom(err error) (val T, found bool) {
	var enrichedErr enrichedErrorWithFields

	for ; err != nil; err = Unwrap(err) {
		if As(err, &enrichedErr) {
			for _, field := range enrichedErr.fields {
				ff, ok := field.(wrappedValue[T])
				if !ok {
					continue
				}
				if ff.typ != f {
					continue
				}
				return ff.val, true
			}
		}
	}

	return val, false
}

func (w wrappedValue[T]) enrich(err *enrichedErrorWithFields) {
	err.fields = append(err.fields, w)
}

type enricher interface {
	enrich(err *enrichedErrorWithFields)
}

type enrichedErrorWithFields struct {
	err    error
	fields []fmt.Stringer
}

func (e enrichedErrorWithFields) Error() string {
	var b strings.Builder
	b.WriteString("original error: ")
	b.WriteString(e.err.Error())
	b.WriteByte('\n')
	b.WriteString("enriched fields:\n")
	for _, f := range e.fields {
		b.WriteString(f.String())
		b.WriteByte('\n')
	}
	return b.String()
}

func (e enrichedErrorWithFields) Unwrap() error {
	return e.err
}

func Enrich(err error, fields ...enricher) error {
	if err == nil {
		return nil
	}
	e := enrichedErrorWithFields{
		err:    err,
		fields: make([]fmt.Stringer, 0, len(fields)),
	}
	for _, field := range fields {
		field.enrich(&e)
	}

	return e
}
