package whoops

import "fmt"

type wrapser interface {
	WrapS(msg string, args ...any) error
}

type wraps struct {
	orig error
	msg  string
}

func WrapS(err error, msg string, args ...any) error {
	return wraps{
		orig: err,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func (w wraps) Error() string {
	return w.orig.Error() + ": " + w.msg
}

func (w wraps) Unwrap() error {
	return w.orig
}
