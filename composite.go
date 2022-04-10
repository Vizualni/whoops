package whoops

type Compositer func(error) error

func Compose(err error, funcs ...Compositer) error {
	for _, fnc := range funcs {
		err = fnc(err)
	}

	return err
}

type customErrorMessage struct {
	err error
	msg string
}

func (c customErrorMessage) Error() string {
	return c.msg
}
func (c customErrorMessage) Unwrap() error {
	return c.err
}

func ComposeCustomMessage(fnc func(error) string) Compositer {
	return func(err error) error {
		return customErrorMessage{
			err: err,
			msg: fnc(err),
		}
	}
}

func ComposeEnrich(fields ...enricher) Compositer {
	return func(err error) error {
		return Enrich(err, fields...)
	}
}
