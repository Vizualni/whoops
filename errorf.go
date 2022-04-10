package whoops

import (
	"fmt"
)

type Errorf string

// Error is fullfilling the error interface only because it's much nicer to
// use it with the Is method.
func (Errorf) Error() string { panic("Errorf is not an error you should return") }

type formattedError struct {
	origErr Errorf
	args    []any
}

func (e formattedError) Is(err error) bool {

	if errf, ok := err.(Errorf); ok {
		return e.origErr == errf
	}
	errf, ok := err.(formattedError)
	if !ok {
		return false
	}

	if e.origErr != errf.origErr {
		return false
	}
	if len(e.args) != len(errf.args) {
		return false
	}

	for i := 0; i < len(e.args); i++ {
		if e.args[i] != errf.args[i] {
			return false
		}
	}

	return true
}

func newFormattedError(origErr Errorf, args ...any) formattedError {
	return formattedError{
		origErr: origErr,
		args:    args,
	}
}

func (e formattedError) Error() string {
	return fmt.Sprintf(string(e.origErr), e.args...)
}

func (e Errorf) Format(args ...any) formattedError {
	return newFormattedError(e, args...)
}
