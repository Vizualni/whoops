package whoops

import (
	"errors"
)

type Field[T any] string // adding T for type safety

type wrappedValue[T any] struct {
	typ Field[T]
	val T
}

func (f Field[T]) Val(v T) wrappedValue[T] {
	return wrappedValue[T]{
		typ: f,
		val: v,
	}
}

func (f Field[T]) GetFrom(err error) (val T, found bool) {
	var we wrapperErrorWithFields

	for ; err != nil; err = errors.Unwrap(err) {
		if errors.As(err, &we) {
			for _, field := range we.fields {
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

func (w wrappedValue[T]) wrap(err *wrapperErrorWithFields) {
	err.fields = append(err.fields, w)
}

type wrapper interface {
	wrap(err *wrapperErrorWithFields)
}

type wrapperErrorWithFields struct {
	err    error
	fields []any
}

func (we wrapperErrorWithFields) Error() string {
	return we.err.Error()
}

func (we wrapperErrorWithFields) Unwrap() error {
	return we.err
}

func Wrap(err error, fields ...wrapper) wrapperErrorWithFields {
	we := wrapperErrorWithFields{
		err:    err,
		fields: make([]any, 0, len(fields)),
	}
	for _, field := range fields {
		field.wrap(&we)
	}

	return we
}
