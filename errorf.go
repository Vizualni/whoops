package whoops

import (
	"fmt"
)

var _ wrapser = formattedErrorf{}

type Errorf string

func (e Errorf) Error() string {
	return fmt.Sprintf("%s - not formatted correctly. Use Format(...) method", string(e))
}

func (e Errorf) Is(err error) bool {
	switch w := err.(type) {
	case Errorf:
		return w == e
	case formattedErrorf:
		return w.Is(e)
	}
	return false
}

func (e Errorf) Format(args ...any) formattedErrorf {
	return formattedErrorf{
		origErr: e,
		msg:     fmt.Sprintf(string(e), args...),
	}
}

type formattedErrorf struct {
	origErr Errorf
	msg     string
}

func (e formattedErrorf) Is(err error) bool {

	if orig, ok := err.(Errorf); ok {
		return e.origErr == orig
	}

	errf, ok := err.(formattedErrorf)
	if !ok {
		return false
	}

	if e.origErr != errf.origErr {
		return false
	}

	if e.msg != errf.msg {
		return false
	}

	return true
}

func (e formattedErrorf) Error() string {
	return e.msg
}

func (e formattedErrorf) WrapS(msg string, args ...any) wraps {
	return WrapS(e, msg, args...)
}
