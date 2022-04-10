package whoops

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
	fields []any
}

func (enrichedErr enrichedErrorWithFields) Error() string {
	return enrichedErr.err.Error()
}

func (enrichedErr enrichedErrorWithFields) Unwrap() error {
	return enrichedErr.err
}

func Enrich(err error, fields ...enricher) enrichedErrorWithFields {
	enrichedErr := enrichedErrorWithFields{
		err:    err,
		fields: make([]any, 0, len(fields)),
	}
	for _, field := range fields {
		field.enrich(&enrichedErr)
	}

	return enrichedErr
}
