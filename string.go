package whoops

type String string

var _ error = String("")
var _ wrapser = String("")

func (s String) Error() string { return (string(s)) }

func (s String) WrapS(msg string, args ...any) error {
	return WrapS(s, msg, args...)
}
