package whoops

import "fmt"

type wrapper struct {
	error
	wrapped error
}

func (w wrapper) Error() string {
	return fmt.Sprintf("base error: %s\nwrapped error: %s", w.error.Error(), w.wrapped.Error())
}

func (w wrapper) Unwrap() error {
	return Group{w.error, w.wrapped}.Unwrap()
}

func Wrap(err, wrap error) wrapper {
	return wrapper{
		error:   err,
		wrapped: wrap,
	}
}
