package whoops

import (
	"fmt"
)

var _ wrapser = formattedError{}

type Errorf string

type isableErrorf struct {
	err Errorf
}

func (e isableErrorf) Error() string {
	panic("not an error you should return. use it with errors.Is")
}

func (e Errorf) CheckIs() error {
	return isableErrorf{e}
}

type formattedError struct {
	origErr Errorf
	args    []any
}

func (e formattedError) Is(err error) bool {

	if isErrf, ok := err.(isableErrorf); ok {
		return e.origErr == isErrf.err
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

func (e Errorf) Format(args ...any) formattedError {
	return newFormattedError(e, args...)
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

func (e formattedError) WrapS(msg string, args ...any) error {
	return WrapS(e, msg, args...)
}
