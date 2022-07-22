package whoops

import "fmt"

type wrapser interface {
	WrapS(msg string, args ...any) wraps
}

type wraps struct {
	orig error
	msg  string
}

func WrapS(err error, msg string, args ...any) wraps {
	return wraps{
		orig: err,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func (w wraps) Error() string {
	return w.orig.Error() + ": " + w.msg
}

func (w wraps) Is(err error) bool {
	w1, ok := err.(wraps)
	if !ok {
		return false
	}

	if w1.msg != w.msg {
		return false
	}

	if Is(w, w1.orig) {
		return true
	}

	return false
}

func (w wraps) Unwrap() error {
	return w.orig
}

func (w wraps) WrapS(msg string, args ...any) error {
	return WrapS(w, msg, args...)
}
