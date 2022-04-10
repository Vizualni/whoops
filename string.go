package whoops

type String string

var _ error = String("")

func (s String) Error() string { return (string(s)) }
